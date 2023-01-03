FROM golang:1.19-bullseye

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY *.go ./
RUN go build -o /rss_service

EXPOSE 8080

ENTRYPOINT ["/rss_service"]
