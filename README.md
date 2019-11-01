## λ Fuzzy Service
Lambda based fuzzy matching function with mutable terms dictionary

#### Installation
Deployment and development of the function requires [serverless](https://github.com/serverless/serverless) & [dep](https://github.com/golang/dep).

```bash
$ npm install -g serverless
```

On Mac, installing **dep** would look something like this.

```bash
$ brew install dep
$ brew upgrade dep
```

#### Build
```bash
$ make build
```

#### Deploy
```bash
$ make deploy
```

#### Links

[Serverless Framework example for Golang and Lambda](https://serverless.com/blog/framework-example-golang-lambda-support/)

[Serverless Examples](https://github.com/serverless/examples)
