package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/maxmclau/fuzzy-service/lib/dictionary"

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
		q := strings.ToUpper(query)

		var match Match

		for _, term := range dict.Terms {
			if strings.Contains(q, strings.ToUpper(term)) {
				match.Terms = append(match.Terms, term)
			}
		}

		match.Query = query

		if len(match.Terms) > 0 {
			matches = append(matches, match)
		}
	}

	if len(matches) > 0 {
		resp, err := json.Marshal(matches)
		if err != nil {
			return utils.ServerError(err)
		}

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(resp),
		}, nil
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNoContent,
		}, nil
	}
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
