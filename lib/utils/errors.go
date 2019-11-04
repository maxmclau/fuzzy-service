package utils

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func ClientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, fmt.Errorf("ClientError: %s", http.StatusText(status))
}

func ServerError(err error) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("ServerError: %s\n", err)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, fmt.Errorf("ServerError: %s", err)
}
