include deploy/.env
include deploy/secret.env
LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

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
	@if [ ! -f "$(LOCAL_BIN)/goose" ]; then \
		echo "Installing goose..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.0; \
	else \
		echo "goose already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/minimock" ]; then \
		echo "Installing minimock..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.5 ; \
	else \
		echo "minimock already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/mockgen" ]; then \
		echo "Installing mockgen..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.2.1; \
	else \
		echo "mockgen already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/grpc-gateway" ]; then \
		echo "Installing grpc-gateway..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.26.3; \
	else \
		echo "grpc-gateway already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/protoc-gen-openapiv2" ]; then \
		echo "Installing protoc-gen-openapiv2..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.26.3; \
	else \
		echo "protoc-gen-openapiv2 already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/statik" ]; then \
		echo "Installing statik..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7; \
	else \
		echo "statik already installed."; \
	fi



get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	mkdir -p pkg/swagger
	make generate-user-api
	make generate-auth-api
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path=api/proto/user_v1 \
		--proto_path=vendor.protogen \
		--go_out=pkg/user_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
		--validate_out lang=go:pkg/user_v1 --validate_opt=paths=source_relative \
        --plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
        --grpc-gateway_out=pkg/user_v1 --grpc-gateway_opt=paths=source_relative \
        --plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
        --openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
        --plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
        api/proto/user_v1/user.proto

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path=api/proto/auth_v1 \
		--proto_path=vendor.protogen \
		--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
		--validate_out=lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
		api/proto/auth_v1/auth.proto

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml


local migration-create:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} create $(name) sql sql

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v


local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v


local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go
copy-to-server:
	scp service_linux root@$(IP_SERVER):

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t $(REGESTRY)/server:v0.0.1 -f deploy/Dockerfile .
	docker login -u $(USERNAME) -p $(PASSWORD) $(REGESTRY)
	docker push $(REGESTRY)/server:v0.0.1

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/Ippolid/auth/internal/service/...,github.com/Ippolid/auth/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/Ippolid/auth/internal/service/...,github.com/Ippolid/auth/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi

grpc-load-test:
	ghz \
        --proto api/proto/user_v1/user.proto \
        -i api/proto,vendor.protogen \
        --call user_v1.UserV1.Get \
        --data '{"id":"28"}' \
        --rps 100 \
        --total 300 \
        --cacert deploy/server_cert.pem \
        localhost:50051


grpc-error-load-test:
	ghz \
		--proto api/proto/user_v1/user.proto \
		-i api/proto,vendor.protogen \
		--call user_v1.UserV1.Get \
		--data '{"id":"-1"}' \
		--rps 100 \
		--total 300 \
		--cacert deploy/server_cert.pem \
		localhost:50051

