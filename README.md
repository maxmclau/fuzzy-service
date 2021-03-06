## λ Fuzzy Service
[![Built for SumUp](https://img.shields.io/badge/Built%20for%20SumUp-blue?style=flat)](http://sumup.com/)

Lambda based fuzzy matching function with mutable terms dictionary

### Installation
Deployment and development of the function requires [serverless](https://github.com/serverless/serverless) & [dep](https://github.com/golang/dep).

```bash
$ npm install -g serverless
```

On Mac, installing **dep** would look something like this.

```bash
$ brew install dep
```

### Build
```bash
$ make build
```

### Deploy
```bash
$ make deploy
$ make deploy-prod
```

### API
Detailed request and response information for API


#### GET /match
Return all matched terms against dictionary
```http
GET /prod/match?q=Ammo&amp; q=I sell ammunition HTTP/1.1
```

<details>
<summary><strong>200 OK</strong> - Response</summary>

```js
[
    {
        "query": "Ammo",
        "terms": [
            "Ammo",
            "Ammunition"
        ]
    }
]
``` 
</details>

```http
GET /prod/match?q=Puppies&amp; q=I sell puppies HTTP/1.1
```

<details>
<summary><strong>204 No Content</strong> - Response</summary>

```js
...
``` 
</details>

#### GET /terms
Return all terms used in fuzzy matching along with the date they were last modified

```http
GET /prod/terms HTTP/1.1
```

<details>
<summary><strong>200 OK</strong> - Response</summary>

```js
{
    "modified": 1572891670,
    "terms": [
        "420",
        "Adult",
        "Airline",
        "Ammo"
        ...
    ]
}
```
</details>

#### POST /terms
Add additional terms to terms dictionary and return updated dictionary

```http
POST /prod/terms HTTP/1.1
Content-Type: application/json
Content-Length: 189

{
    "terms": [
        "Coffee",
        "Theft"
    ]
}
```

<details>
<summary><strong>200 OK</strong> - Response</summary>

```js
{
    "modified": 1572893178,
    "terms": [
        "420",
        "Adult",
        "Airline",
        "Ammo",
        ...
        "Coffee",
        "Theft"
    ]
}
```
</details>

#### PUT /terms
Replace all terms in dictionary and return updated dictionary

```http
POST /prod/terms HTTP/1.1
Content-Type: application/json
Content-Length: 167

{
    "terms": [
        "Island",
        "Epstein"
    ]
}
```

<details>
<summary><strong>200 OK</strong> - Response</summary>

```js
{
    "modified": 1572893446,
    "terms": [
        "Island",
        "Epstein"
    ]
}
``` 
</details>
    
### Links

[Serverless Framework example for Golang and Lambda](https://serverless.com/blog/framework-example-golang-lambda-support/)

[Serverless Examples](https://github.com/serverless/examples)
