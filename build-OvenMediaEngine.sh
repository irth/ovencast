#!/usr/bin/env bash

DOCKER_ORG="${DOCKER_ORG:-irth7}"
version="$(
    cd OvenMediaEngine
    git describe --tags
)"

docker buildx build \
    --push \
    --platform linux/arm/v7,linux/arm64/v8,linux/amd64 \
    --tag "$DOCKER_ORG/ovenmediaengine:$version" \
    ./OvenMediaEngine
