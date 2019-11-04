package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/maxmclau/fuzzy-service/lib/dictionary"

	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/maxmclau/fuzzy-service/lib/utils"
)

var (
	AWS_BUCKET = os.Getenv("AWS_BUCKET")
)

type Match struct {
	Query string   `json:"query"`
	Terms []string `json:"terms"`
}

func Get(sess *session.Session, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queries := req.MultiValueQueryStringParameters["q"]
	if len(queries) == 0 {
		return utils.ClientError(http.StatusBadRequest)
	}

	var dict dictionary.Dictionary

	err := dictionary.GetDictionary(sess, &dict, AWS_BUCKET, false)
	if err != nil {
		return utils.ServerError(err)
	}

	var matches []Match

	for _, query := range queries {
		terms := fuzzy.FindFold(query, dict.Terms)

		if len(terms) != 0 {
			var match Match

			match.Query = query
			match.Terms = terms

			matches = append(matches, match)
		}
	}

	resp, err := json.Marshal(matches)
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
		return Get(sess, req)
	default:
		return utils.ClientError(http.StatusMethodNotAllowed)
	}
}

func main() {
	lambda.Start(Handler)
}
