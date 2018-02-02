#!/bin/bash

set -euxo pipefail

docker run --rm -t \
    -p 10000:10000 \
    -p 10001:10001 \
    -p 8080:8080 \
    {{.ServiceName}}
