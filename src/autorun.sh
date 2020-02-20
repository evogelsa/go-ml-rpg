#!/bin/sh
echo 'Retrieving dependencies'
go get
echo 'Running server on localhost:8080/'
go run main.go
