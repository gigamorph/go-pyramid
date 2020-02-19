#!/bin/sh

BUILD_FLAGS="-p 1"
TEST_FLAGS="-count 1 -v"

set -x
go test $BUILD_FLAGS $TEST_FLAGS ./...
