GO_PATH := $(shell go env GOPATH)

GO_PACKAGE := github.com/kazhuravlev/srv-email
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*.pb.go" -not -path "*.pb.gw.go")
PROTO_TARGET := ./contracts

.PHONY: format
format: format-code

.PHONY: format-code
format-code:
	@echo 'format-code'

	go install github.com/Eun/goremovelines/cmd/goremovelines
	go install golang.org/x/tools/cmd/goimports

	@goremovelines -w $(GO_FILES)
	@goimports -local $(GO_PACKAGE) -w $(GO_FILES)
	@gofmt -w $(GO_FILES)

.PHONY: generate
generate: generate-code format

.PHONY: generate-code
generate-code:
	@echo 'generate code'

	go install github.com/kazhuravlev/options-gen/cmd/options-gen

	@which protoc-gen-go >/dev/null || GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go

	mkdir -p $(PROTO_TARGET)
	protoc \
		-I ./contracts \
		-I $(GO_PATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.1 \
		-I $(GO_PATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.1/third_party/googleapis/ \
		-I $(GO_PATH)/src \
		--go_out=plugins=grpc:$(PROTO_TARGET) \
		--grpc-gateway_out=logtostderr=true:$(PROTO_TARGET) \
		--swagger_out=logtostderr=true:$(PROTO_TARGET) \
		./contracts/*.proto

	go generate ./...
