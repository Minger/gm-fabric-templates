#!/bin/bash

set -euxo pipefail

# assume we are in the service directory
SERVICEDIR=$PWD
DOCKERDIR="${SERVICEDIR}/docker"

(
	cd "cmd/server"
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o "$DOCKERDIR/{{.ServiceName}}" .
)

(
	cp "settings.toml" "${DOCKERDIR}/."
	cd $DOCKERDIR
	docker build -t {{.ServiceName}} .
)

