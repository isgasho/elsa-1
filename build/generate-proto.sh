#!/usr/bin/env bash
echo 'start generate proto...'


echo "cd ../pkg/proto"
cd ../pkg/proto
protoc --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative ./registry.proto
echo "generate the pb files success"