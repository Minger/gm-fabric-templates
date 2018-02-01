#!/bin/bash

set -euxo pipefail

docker run --rm -t \
    -p {{.GrpcServerPort}}:{{.GrpcServerPort}} \
    -p {{.MetricsServerPort}}:{{.MetricsServerPort}} \
    -p {{.GatewayProxyPort}}:{{.GatewayProxyPort}} \
    {{.ServiceName}}
