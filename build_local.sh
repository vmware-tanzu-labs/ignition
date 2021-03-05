#!/bin/bash

set -e

PROJECT_DIR=$(dirname $0)
[[ $PROJECT_DIR == "." ]] && PROJECT_DIR=$(pwd)

BUILD_DIR="$PROJECT_DIR/build"

mkdir -p $BUILD_DIR

pushd $PROJECT_DIR
  go test ./...

  pushd $PROJECT_DIR/cmd/ignition
    GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/ignition-linux
    GOOS=darwin GOARCH=amd64 go build -o $BUILD_DIR/ignition-mac
    GOOS=windows GOARCH=amd64 go build -o $BUILD_DIR/ignition-win64
  popd

  pushd $PROJECT_DIR/web
    yarn install && yarn lint && yarn test && yarn build
    cp -r dist/* $BUILD_DIR/
  popd

popd
