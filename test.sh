#!/bin/bash

set -e

mockery
go test -tags=integration ./...