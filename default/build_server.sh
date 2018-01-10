#!/bin/bash

set -euxo pipefail

# assume we are in the service directory
SERVICEDIR=$PWD
SERVERDIR="${SERVICEDIR}/cmd/server"

(
	cd $SERVERDIR
	go build -o=$GOPATH/bin/{{.ServiceName}}
)
