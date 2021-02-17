#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  greet/greetpb/greet.proto

protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  sum/sumpb/sum.proto

protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  prime/primepb/prime.proto

# protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
# protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
# protoc blog/blogpb/blog.proto --go_out=plugins=grpc:.
