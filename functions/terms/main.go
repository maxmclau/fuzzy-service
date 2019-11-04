package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/maxmclau/fuzzy-service/lib/dictionary"
	"github.com/maxmclau/fuzzy-service/lib/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	AWS_BUCKET = os.Getenv("AWS_BUCKET")
)

func Get(sess *session.Session) (events.APIGatewayProxyResponse, error) {
	var dict dictionary.Dictionary

	err := dictionary.GetDictionary(sess, &dict, AWS_BUCKET, false)
	if err != nil {
		return utils.ServerError(err)
	}

	resp, err := json.Marshal(dict)
	if err != nil {
		return utils.ServerError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(resp),
	}, nil
}

func PostPut(sess *session.Session, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return utils.ClientError(http.StatusNotAcceptable)
	}

	var requestDict dictionary.Dictionary
	err := json.Unmarshal([]byte(req.Body), &requestDict)
	if err != nil {
		return utils.ClientError(http.StatusUnprocessableEntity)
	}

	if len(requestDict.Terms) == 0 {
		return utils.ClientError(http.StatusBadRequest)
	}

	if req.HTTPMethod == http.MethodPost {
		var currentDict dictionary.Dictionary
		err = dictionary.GetDictionary(sess, &currentDict, AWS_BUCKET, false)
		if err != nil {
			return utils.ServerError(err)
		}

		requestDict.Terms = append(requestDict.Terms, currentDict.Terms...)
	}

	requestDict.Modified = time.Now().Unix()

	resp, err := dictionary.SetDictionary(sess, &requestDict, AWS_BUCKET)
	if err != nil {
		return utils.ServerError(err)
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
		return utils.ClientError(http.StatusMethodNotAllowed)
	}
}

func main() {
	lambda.Start(Handler)
}
