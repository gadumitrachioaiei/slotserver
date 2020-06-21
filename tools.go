// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-check-breaking"
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-check-lint"
	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc" // we need this to use the new plugin for grpc
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
