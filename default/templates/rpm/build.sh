#!/bin/sh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SRC="$( cd "$DIR/.." && pwd )"

IMAGE_ID=$(docker build -q $SRC/pkg | cut -d: -f 2 -)
IMAGE_ID=${IMAGE_ID:0:12}

docker run --rm -i -v $SRC:/src "$IMAGE_ID" /bin/bash -s < $SRC/pkg/script
