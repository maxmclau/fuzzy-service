package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetSSMParameter(sess *session.Session, name string, decrypt bool) (string, error) {
	// create an instance of the SSM Session
	ssmSession := ssm.New(sess)

	// create the request to SSM
	getParameterInput := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(decrypt),
	}

	// cet the parameter from SSM
	param, err := ssmSession.GetParameter(getParameterInput)
	if err != nil {
		return "", err
	}

	return *param.Parameter.Value, nil
}
