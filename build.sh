#!/usr/bin/env bash
set -e

version=`git rev-parse HEAD`
last_build=`date +"%Y/%m/%d-%H:%M:%S"`
go build -ldflags "-X main.Version=${version} -X main.LastBuild=${last_build}" imghr.go
