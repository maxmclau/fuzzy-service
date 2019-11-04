package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	sess *session.Session
}

func (s *MatchTestSuite) SetupSuite() {
	awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

	s.sess = session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: awsCredentials,
	}))
}

func (s *MatchTestSuite) TearDownSuite() {
	os.Remove(filepath.Join(testDataDir, testDataDownloadFilename))

	s3Session := s3.New(s.sess)
	s3Session.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(awsBucket), Key: aws.String(testDataUploadFilename)})
}

func (s *MatchTestSuite) TestDownloadExistingFile() {
	err := DownloadFromS3(s.sess, testDataDir, testDataDownloadFilename, awsBucket)

	assert.NoError(s.T(), err)
}

func (s *MatchTestSuite) TestDownloadNotExistingFile() {
	err := DownloadFromS3(s.sess, testDataDir, "epstein.txt", awsBucket)

	assert.Error(s.T(), err)
}

func (s *MatchTestSuite) TestUploadNotExistingFile() {
	err := UploadToS3(s.sess, testDataDir, "epstein.txt", awsBucket)

	assert.Error(s.T(), err)
}

func (s *MatchTestSuite) TestUploadExistingFile() {
	err := UploadToS3(s.sess, testDataDir, testDataUploadFilename, awsBucket)

	assert.NoError(s.T(), err)
}

func TestMatchTestSuite(t *testing.T) {
	suite.Run(t, new(MatchTestSuite))
}
