#!/usr/bin/env bash
set -e

go build -ldflags "-X main.Version=`git rev-parse HEAD`" imghr.go
