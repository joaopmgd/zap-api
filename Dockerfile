# Based on the latest golang image from Docker
FROM golang:latest AS builder

ADD . /go/src/gitlab.com/api
WORKDIR /go/src/gitlab.com/api

ENV ZAP_PROPERTIES_ENDPOINT "http://grupozap-code-challenge.s3-website-us-east-1.amazonaws.com/sources/source-2.json"

# Testing the API
RUN go test

# Building the Go executable for linux
RUN env GOOS=linux GOARCH=386 go build -o /app -i main.go

########################################################

FROM scratch

ENV ZAP_PROPERTIES_ENDPOINT "http://grupozap-code-challenge.s3-website-us-east-1.amazonaws.com/sources/source-2.json"
COPY --from=builder /app ./
WORKDIR /

ENV HOST=':8080'
EXPOSE 8080

ENTRYPOINT ["./app"]