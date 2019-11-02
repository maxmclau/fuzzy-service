package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Dictionary struct {
	Modified int64    `json:"modified,omitempty"`
	Terms    []string `json:"terms"`
}

const (
	DictionaryBucket        = "fuzzy-service-bucket"
	DictionaryFile          = "dictionary.json"
	DictionaryFileDirectory = "/tmp"
	DictionaryLifespan      = 1800
)

var path = filepath.Join(DictionaryFileDirectory, DictionaryFile)

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func getDictionary(session *session.Session, dictionary *Dictionary, forceUpdate bool) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		// download from s3 if file doesn't exist
		err = DownloadFromS3(session, DictionaryFileDirectory, DictionaryFile, DictionaryBucket)
	} else if err != nil {
		// return server error if unrecognized error
		return err
	} else {
		// check age of file
		now := time.Now()
		age := now.Sub(info.ModTime()).Seconds()

		// a busy Lamdba function could go a while before exiting so we want to be sure it
		// periodically pulls new versions of our dictionary
		if age > DictionaryLifespan || forceUpdate {
			err = DownloadFromS3(session, DictionaryFileDirectory, DictionaryFile, DictionaryBucket)
		}
	}

	if err != nil {
		return err
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into dictionary
	err = json.Unmarshal(byteValue, &dictionary)
	if err != nil {
		return err
	}

	return nil
}

func Get(session *session.Session) (events.APIGatewayProxyResponse, error) {
	var dictionary Dictionary

	err := getDictionary(session, &dictionary, false)
	if err != nil {
		return serverError(err)
	}

	resp, err := json.Marshal(dictionary)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(resp),
	}, nil
}

func PostPut(session *session.Session, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	var requestDictionary Dictionary
	err := json.Unmarshal([]byte(req.Body), &requestDictionary)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	if len(requestDictionary.Terms) == 0 {
		return clientError(http.StatusBadRequest)
	}

	if req.HTTPMethod == http.MethodPost {
		var dictionary Dictionary
		err = getDictionary(session, &dictionary, true)
		if err != nil {
			return serverError(err)
		}

		requestDictionary.Terms = append(requestDictionary.Terms, dictionary.Terms...)
	}

	requestDictionary.Modified = time.Now().Unix()

	resp, err := json.Marshal(requestDictionary)
	if err != nil {
		return serverError(err)
	}

	err = ioutil.WriteFile(path, []byte(resp), 0666)
	if err != nil {
		return serverError(err)
	}

	err = UploadToS3(session, DictionaryFileDirectory, DictionaryFile, DictionaryBucket)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(resp),
	}, nil
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())

	switch req.HTTPMethod {
	case http.MethodGet:
		return Get(sess)
	case http.MethodPost, http.MethodPut:
		return PostPut(sess, req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func main() {
	lambda.Start(Handler)
}
