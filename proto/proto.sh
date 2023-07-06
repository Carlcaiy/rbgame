#!/bin/sh
protoc --proto_path=. --go_out=. login.proto
protoc --proto_path=. --go_out=. poker.proto
protoc --proto_path=. --go_out=. test.proto

protoc --go_out=. --go-grpc_out=. greet.proto
