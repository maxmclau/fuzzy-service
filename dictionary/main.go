package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var awsSess *session.Session
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

const (
	DictionaryStore = "config-restricted-dictionary"
)

type ResponseBody struct {
	Dictionary string `json:"dictionary"`
}

type RequestBody struct {
	Terms []string `json:"terms"`
}

func Post(sess *session.Session, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	body := new(RequestBody)

	err := json.Unmarshal([]byte(req.Body), body)

	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	err = addTerms(sess, *body)

	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Location": "/dictionary"},
	}, nil
}

func addTerms(sess *session.Session, terms RequestBody) error {
	svc := ssm.New(sess)

	dictionary, err := getDictionary(sess)

	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	buf.WriteString(dictionary)
	buf.WriteString(term)
	appended := buf.String()

	_, err = svc.PutParameter(
		&ssm.PutParameterInput{
			Name:      aws.String(DictionaryStore),
			Value:     aws.String(appended),
			Type:      aws.String("StringList"),
			Overwrite: aws.Bool(true),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func getDictionary(sess *session.Session) (string, error) {
	svc := ssm.New(sess)

	output, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name:           aws.String(DictionaryStore),
			WithDecryption: aws.Bool(false),
		},
	)

	return *output.Parameter.Value, err
}

func Get(sess *session.Session, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dictionary, err := getDictionary(sess)

	if err != nil {
		return serverError(err)
	}

	resp, _ := json.Marshal(ResponseBody{Dictionary: dictionary})

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(resp),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if awsSess == nil {
		awsSess = session.Must(session.NewSession())
	}

	switch req.HTTPMethod {
	case "GET":
		return Get(awsSess, req)
	case "POST":
		return Post(awsSess, req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func main() {
	lambda.Start(handler)
}
