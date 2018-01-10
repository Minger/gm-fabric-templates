#!/bin/bash

date="$(date -u '+%Y-%m-%d%-H%-M%-S')"

version_major=1
version_minor=0
version_patch=0
version="$version_major.$version_minor.$version_patch"

PATH=/usr/local/go/bin:$PATH:$(pwd)/.go/bin
export GOPATH=$(pwd)/.go
mkdir -p .{{.ServicePath}}
ln -s /src .{{.ServicePath}}
cd .{{.ServicePath}}

prefix=build/opt/services/{{.ServiceName}}-${version_major}.${version_minor}
mkdir -p $prefix/bin
mkdir -p $prefix/docs

# Copy stuff over
go build -v -o $prefix/bin/{{.ServiceName}}-bin ./cmd/server

sed < pkg/init.py > $prefix/bin/{{.ServiceName}} \
    -e "s/{{.ServiceName}}-X\\.X/{{.ServiceName}}-${version_major}.${version_minor}/g" \
    -e "s/^VERSION_MAJOR =.*/VERSION_MAJOR = \"${version_major}\"/" \
    -e "s/^VERSION_MINOR =.*/VERSION_MINOR = \"${version_minor}\"/" \
    -e "s/^VERSION_PATCH =.*/VERSION_PATCH = \"${version_patch}\"/"
    chmod +x $prefix/bin/{{.ServiceName}}


// cp -r docs/static/* $prefix/docs

# Build the RPM
rm *.rpm
fpm -C build -s dir -t rpm -n {{.ServiceName}} -v "${version}_RC-${date}" opt
