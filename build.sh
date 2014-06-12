#!/bin/bash

export GOPATH=`pwd`
export GOBIN="${GOPATH}/bin"

if [ ! -d "src/github.com/garyburd/redigo/redis" ]; then
  go get github.com/garyburd/redigo/redis
fi

if [ -d "pkg" ]; then
  rm -Rf pkg/*
fi

go install check_redis.go
