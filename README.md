# Golang RSS service

## About the project

The service goal is to parse and return RSS feed as fast as possible. The service implements JWT token authorization.

### Prerequisites

Go version go1.19.2 or higher

### Development

#### Download the repo

```
git clone https://github.com/Alexander96/rss-service
```

#### Install dependencies

```
go mod download
```

#### Running it in development mode

```
go run main.go
```

#### Building a Docker image and running it

```
docker build -t rss-service .
docker run -d -p 808:8080 -t rss-service
```

## Usage

First generate a token. Currently the service supports only admin admin acc

```
curl -i -u 'admin:admin' --request GET 'localhost:8080/login'
```

Then make GET API call using the token in a "token" header and body array of RSS URLs

```
curl --request GET 'localhost:8080/rss' \
--header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNjcyNzc0NDc2fQ.y6BPWEvKp3Ow7mEoBh-BYjB4nvYYV7PqGPZsHZbBeXQ' \
--header 'Content-Type: application/json' \
--data-raw '{
    "urls": [
        "http://rss.cnn.com/rss/edition.rss",
        "http://rss.cnn.com/rss/edition_world.rss",
        "http://rss.cnn.com/rss/edition_asia.rss",
        "http://rss.cnn.com/rss/edition_travel.rss",
        "http://rss.cnn.com/rss/cnn_latest.rss"
    ]
}'
```
