package dictionary

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/maxmclau/fuzzy-service/lib/utils"
)

const (
	lifespan  float64 = 1800
	directory string  = "/tmp"
	filename  string  = "dictionary.json"
)

var path = filepath.Join(directory, filename)

type Dictionary struct {
	Modified int64    `json:"modified,omitempty"`
	Terms    []string `json:"terms"`
}

func SetDictionary(sess *session.Session, dict *Dictionary, bucket string) (string, error) {
	dictJson, err := json.Marshal(dict)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(path, []byte(dictJson), 0666)
	if err != nil {
		return "", err
	}

	err = utils.UploadToS3(sess, directory, filename, bucket)
	if err != nil {
		return "", err
	}

	return string(dictJson), nil
}

func GetDictionary(sess *session.Session, dict *Dictionary, bucket string, forceUpdate bool) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		// download from s3 if file doesn't exist
		err = utils.DownloadFromS3(sess, directory, filename, bucket)
	} else if err != nil {
		// return server error if unrecognized error
		return err
	} else {
		// check age of file
		now := time.Now()
		age := now.Sub(info.ModTime()).Seconds()

		// a busy Lamdba function could go a while before exiting so we want to be sure it
		// periodically pulls new versions of our dictionary
		if age > lifespan || forceUpdate {
			err = utils.DownloadFromS3(sess, directory, filename, bucket)
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
	err = json.Unmarshal(byteValue, &dict)
	if err != nil {
		return err
	}

	return nil
}
