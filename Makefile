.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/terms ./functions/terms
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/match ./functions/match

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

deploy-prod: clean build
	sls deploy -s prod --verbose