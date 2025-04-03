LOCAL_BIN:=$(CURDIR)/bin

#LOCAL_BIN:=$(CURDIR)/bin
#LOCAL_BIN:=$(CURDIR)/../bin

install-deps:
	@if [ ! -f "$(LOCAL_BIN)/protoc-gen-go" ]; then \
		echo "Installing protoc-gen-go..."; \
		GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1; \
	else \
		echo "protoc-gen-go already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/protoc-gen-go-grpc" ]; then \
		echo "Installing protoc-gen-go-grpc..."; \
		GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2; \
	else \
		echo "protoc-gen-go-grpc already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/golangci-lint" ]; then \
		echo "Installing golangci-lint..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0; \
	else \
		echo "golangci-lint already installed."; \
	fi

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-auth-api

generate-auth-api:
	mkdir -p grpc/pkg/auth_v1
	protoc --proto_path grpc/api/auth_v1 \
	--go_out=grpc/pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=grpc/pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	grpc/api/auth_v1/auth.proto

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml