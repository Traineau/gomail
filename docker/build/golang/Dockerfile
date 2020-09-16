# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.12-alpine base image
FROM golang:1.14.1-alpine3.11 as builder

ENV GO111MODULE=on

# The latest alpine image don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Add Maintainer Info
LABEL maintainer="Equipe 01 - Hetic"

WORKDIR /go/src/gomail
COPY go.sum go.mod /go/src/gomail/

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . /go/src/gomail/

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/gomail .

FROM alpine

COPY --from=builder /go/src/gomail/bin/gomail /

EXPOSE 8080

ENTRYPOINT ["/gomail"]