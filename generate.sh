#!/bin/bash

MODULE="github.com/getsynq/entity-status-api/protos"
GEN_DIR="gen/golang"

rm -rf $GEN_DIR

mkdir -p $GEN_DIR

protoc --proto_path=./protos \
    -I ./protos \
    -I /include \
    --go_out=./$GEN_DIR \
    --go_opt=module=$MODULE \
    --go-grpc_out=./$GEN_DIR \
    --go-grpc_opt=module=$MODULE \
    ./protos/run.proto ./protos/entity_status_service.proto

echo "module $MODULE

go 1.23

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.0-20240401165935-b983156c5e99.1
	go.uber.org/mock v0.4.0
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.34.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240610135401-a8a62080eff3 // indirect
)" > gen/golang/go.mod
