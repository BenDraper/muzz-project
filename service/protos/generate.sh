#!/usr/bin/env bash

rm ./*.pb.go

protoc \
  --proto_path=. \
  *.proto \
  --go_out=. \
  --go_opt=paths=source_relative explore-service.proto \
  --go-grpc_out=. \
  --go-grpc_opt=require_unimplemented_servers=false \
  --go-grpc_opt=paths=source_relative explore-service.proto

mockgen -source=../protos/explore-service_grpc.pb.go -destination=../mocks/explore-service.go -package=mocks
