#!/usr/bin/env bash
set -e

go build -ldflags "-X main.Version `git rev-parse HEAD` -X main.LastBuild \"`date +"%Y/%m/%d %H:%M:%S"`\"" imghr.go
