#!/bin/bash

set -e

go test -tags=integration -count=1 ./...