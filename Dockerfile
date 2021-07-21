FROM golang:1.16.5-alpine3.14 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 

#install git
RUN apk update && apk add --no-cache git

# Set current working directory
WORKDIR /app

COPY go.mod .

COPY go.sum . 

RUN go mod download

COPY . . 

RUN go build -o goauth ./cmd/api

FROM alpine:3.13 

WORKDIR /app

COPY  --from=builder  /app/goauth .

EXPOSE 8000

CMD ["/app/goauth"]