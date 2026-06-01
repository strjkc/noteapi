FROM golang:1.26.1-alpine as builder
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR ./app
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY auth ./auth
COPY converter ./converter
COPY database ./database
COPY handlers ./handlers
COPY internal ./internal
COPY spellcheck ./spellcheck
COPY state ./state
COPY storage ./storage
COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY sqlc.yaml ./sqlc.yaml
COPY main.go ./main.go
RUN sqlite3 notes.db
RUN go build -o server .
EXPOSE 8082
