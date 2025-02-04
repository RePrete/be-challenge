# syntax=docker/dockerfile:1

FROM golang:1.23

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      unzip

RUN curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-linux-x86_64.zip" \
    && unzip protoc-25.1-linux-x86_64.zip -d / \
    && protoc --version

RUN go env GOPATH

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 \
    && export PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /app

ADD protos protos

COPY generate.sh ./

RUN chmod +x ./generate.sh \
    && ./generate.sh

COPY go.mod go.sum ./

RUN go mod download

ADD ./app ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

CMD ["/server"]
