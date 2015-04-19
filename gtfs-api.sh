#!/bin/bash

export GTFS_DB_PORT=5432
export GTFS_DB_HOSTNAME=127.0.0.1
export GTFS_DB_DIALECT=postgres
export GTFS_DB_DATABASE=gtfs
export GTFS_DB_USERNAME=gtfs
export GTFS_DB_PASSWORD=gtfs
export REDIS_PORT=8888
export HTTP_PORT=4000
#export GTFS_DB_MIN_CNX=128
#export GTFS_DB_MAX_CNX=128

go clean
go build
go run app.go
