#!/bin/sh

set -e

BUILD_TYPE="dev"
BUILD_VERSION=$(git rev-list -1 --abbrev-commit HEAD)
BUILD_DATE="$(date -u)"
IMAGE_DATE="$(date +%F)"

# Build Docker image for the local arch
build_docker_image() {
  docker build --build-arg BUILD_DATE="$BUILD_DATE" \
                --build-arg BUILD_TYPE="$BUILD_TYPE" \
                --build-arg BUILD_VERSION="$BUILD_VERSION" \
                -t order-book:dev \
                -f deploy/Dockerfile .
}

# Build binary application to run it in local
local_go_build() {
  go mod tidy
  go build \
      -tags debug \
      -ldflags "-X 'main.BuildDateTime=$BUILD_DATE' \
                -X 'main.BuildVersion=$BUILD_VERSION' \
                -X 'main.BuildType=$BUILD_TYPE'" \
      -o out/order-book \
      github.com/lkolacz/OrderBook/cmd/service
}

if [ -n "$1" ]; then
  case $1 in
    docker)
      build_docker_image
      ;;
  esac
else
  local_go_build
fi
