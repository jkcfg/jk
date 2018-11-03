#!/bin/sh -e

# Build jk
export GO111MODULE=on
go install

# Tests, both unit tests and integration tests under /tests
go test -v ./...
