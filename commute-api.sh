#!/bin/bash

export HTTP_PORT=4000

export DB_PORT=5432
export DB_HOSTNAME=127.0.0.1
export DB_DIALECT=postgres
export DB_DATABASE=commute
export DB_USERNAME=commute
export DB_PASSWORD=commute
#export DB_MIN_CNX=128
#export DB_MAX_CNX=128

export REDIS_PORT=8888

go clean
go build
go run app.go
