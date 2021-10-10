#!/bin/bash

##protoc --go_out=. --go_opt=paths=source_relative greet/greetpb/greet.proto

protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative greet/greetpb/greet.proto
protoc  calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.  --go_opt=paths=source_relative
##protoc --go_out=plugins=grpc:. *.proto