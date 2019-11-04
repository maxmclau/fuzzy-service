package utils

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func DownloadFromS3(sess *session.Session, directory string, filename string, bucket string) error {
	path := filepath.Join(directory, filename)

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	s3Downloader := s3manager.NewDownloader(sess)

	_, err = s3Downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		os.Remove(path)

		return err
	}

	return nil
}

func UploadToS3(sess *session.Session, directory string, filename string, bucket string) error {
	path := filepath.Join(directory, filename)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)

	// read file into buffer
	file.Read(buffer)

	s3Uploader := s3.New(sess)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3Uploader.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(filename),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
