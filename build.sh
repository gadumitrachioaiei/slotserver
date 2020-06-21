#!/bin/bash

set -xeuo pipefail
IFS=$'\n\t'
shopt -s expand_aliases

# install protoc
version=$(protoc --version || echo "missing")
if [[ "$version" != "libprotoc 3.12.1" ]]; then
  PB_REL="https://github.com/protocolbuffers/protobuf/releases"
  curl -L $PB_REL/download/v3.12.1/protoc-3.12.1-linux-x86_64.zip -o  /tmp/protoc-3.12.1.zip
  unzip /tmp/protoc-3.12.1.zip -d /tmp/protoc
  cp /tmp/protoc/bin/protoc "$HOME"/.local/bin/
  alias protoc="$HOME"/.local/bin/
fi

# we run the proto linter
#buf check lint

# compile the protocol buffers files
MODULE=github.com/gadumitrachioaiei/slotserver

rm -rf proto/gen/go/*
rm -rf proto/gen/swagger/*

files=$(find proto -name '*.proto' -print0 | xargs -0 -L1 echo)
protoc -I third_party/proto:proto --go_out . --go_opt=module="$MODULE" ${files}
protoc -I third_party/proto:proto --go-grpc_out=module=$MODULE:. ${files}

# generate the http api that the grpc gateway proxy exposes, together with its documentation
for dir in $(echo "${files[@]}" | xargs dirname | sort | uniq); do
  protoc -I third_party/proto:proto --grpc-gateway_out=paths=source_relative:proto/gen/go "${dir}"/*.proto
  protoc -I third_party/proto:proto --openapiv2_out docs/gen "${dir}"/*.proto
done

# generate swagger documentation for our http api, from specially written and annotated code in docs folder
swagger generate spec -w docs -o swagger.json && (swagger validate swagger.json 2>/dev/null || (rm swagger.json && exit 1))
