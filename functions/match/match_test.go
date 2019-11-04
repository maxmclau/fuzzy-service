package main

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	awsRegion          = os.Getenv("AWSREGION")
	awsBucket          = os.Getenv("AWSBUCKET")
	awsAccessKeyID     = os.Getenv("AWSACCESSKEYID")
	awsSecretAccessKey = os.Getenv("AWSSECRETACCESSKEY")
)

const (
	testDataDir              = "./testdata"
	testDataUploadFilename   = "test-s3-upload.json"
	testDataDownloadFilename = "test-s3-download.json"
)

type MatchTestSuite struct {
	suite.Suite
	sess    *session.Session
	testDir string
}

func (suite *MatchTestSuite) SetupSuite() {
	awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

	suite.sess = session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: awsCredentials,
	}))
}

/*
func (suite *MatchTestSuite) TearDownSuite() {
	s3Session := s3.New(suite.sess)
	s3Session.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(awsBucket), Key: aws.String("helloworld.txt")})
}*/

func (suite *MatchTestSuite) TestGetWithoutParameters() {
	_, err := Get(suite.sess, events.APIGatewayProxyRequest{})

	assert.Error(suite.T(), err)
}

func TestMatchTestSuite(t *testing.T) {
	suite.Run(t, new(MatchTestSuite))
}
