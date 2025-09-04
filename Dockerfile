# syntax=docker/dockerfile:1

FROM golang:1.23 AS build

WORKDIR $GOPATH/src/github.com/brotherlogic/kubebrainz

COPY go.mod ./
#COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /kubebrainz

##
## Deploy
##
FROM ubuntu:22.04
USER root:root

# Install ca-certificates package
RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /
COPY --from=build /kubebrainz /kubebrainz

EXPOSE 8081

ENTRYPOINT ["/kubebrainz"]