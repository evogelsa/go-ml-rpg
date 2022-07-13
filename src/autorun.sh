#!/bin/sh
echo 'Retrieving dependencies'
go get
echo 'Running server on localhost:8081/'
go run main.go
