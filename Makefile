
build:
	go build -o main cmd/server/main.go

run: build
	./main

grpc:
	protoc --go_out=. --go-grpc_out=. proto/executor/v1/executor.proto proto/problems/v1/problems.proto
