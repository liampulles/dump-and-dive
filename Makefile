PROTOC_VERSION=3.14.0
REMOTE=github.com/liampulles/dump-and-dive

# --- Common ---
coverage.txt:
	go test -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./pkg/...,./ ./...
view-cover: clean coverage.txt
	go tool cover -html=coverage.txt
test: build
	go test ./test/...
build:
	go build ./...
install: build
	go install ./...
inspect: build
	golint ./...
update:
	go get -u ./...
pre-commit: update grpc-gen clean coverage.txt inspect
	go mod tidy
clean:
	rm -f coverage.txt

# --- GRPC Related ---
grpc-gen: $(GOBIN)/protoc $(GOBIN)/protoc-gen-go $(GOBIN)/protoc-gen-go-grpc $(GOBIN)/protoc-gen-grpc-gateway $(GOBIN)/protoc-gen-openapiv2 tmp/googleapis
	protoc \
   -I ./api/proto \
   -I tmp/googleapis \
   -I tmp/protoc/include \
   -I /usr/local/include \
   -I $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.3.0/ \
   --go_out . \
   --go_opt module=$(REMOTE) \
   --go-grpc_out . \
   --go-grpc_opt module=$(REMOTE) \
   --grpc-gateway_out . \
   --grpc-gateway_opt module=$(REMOTE) \
   --grpc-gateway_opt logtostderr=true \
   --openapiv2_out ./api/openapi/ \
   api/proto/todo_service/todo_service.proto

$(GOBIN)/protoc: tmp/protoc/bin/protoc
	cp tmp/protoc/bin/protoc $(GOBIN)
$(GOBIN)/protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
$(GOBIN)/protoc-gen-go-grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
$(GOBIN)/protoc-gen-grpc-gateway:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
$(GOBIN)/protoc-gen-openapiv2:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
tmp/protoc/bin/protoc:
	wget -O tmp/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
	unzip tmp/protoc.zip -d tmp/protoc
tmp/googleapis:
	wget -O tmp/googleapis.zip https://github.com/googleapis/googleapis/archive/master.zip
	unzip tmp/googleapis.zip -d tmp
	mv tmp/googleapis-master tmp/googleapis
	rm tmp/googleapis.zip