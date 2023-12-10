#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

cd $SCRIPT_DIR/../..

go build

cd examples/simple

../../bardock service1
../../bardock service2
../../bardock service1 service2
