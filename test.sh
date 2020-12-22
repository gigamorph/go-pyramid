#!/bin/sh

BUILD_FLAGS="-p 1"
#TEST_FLAGS="-count 1 -v"
TEST_FLAGS="-count 1"

set -x
go test $BUILD_FLAGS $TEST_FLAGS ./...
