PLATFORM_INFO:=$$(uname)-x$$(getconf LONG_BIT)
PROTO_ENTRY:=op-service.proto
PROTO_ENTRY_FOLDER:=api/proto/v1
LAST_COMMIT:=$$(git rev-list -1 HEAD)
TAG:=$$(git describe)
GO_VERSION:=$$(go version)

all: proto restgateway swagger generatemodel releaseb

generatemodel:
	go generate ./pkg/db

installdepend:
	go get -u -v github.com/golang/protobuf/protoc-gen-go;
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway;
	go get -u -v github.com/golang/protobuf/protoc-gen-go;


proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@ echo "Compiling protobufs to go definitions..."
	@ for file in $(shell ls $(PROTO_ENTRY_FOLDER)/*.proto); do \
		echo "proto - $$file"; \
		protoc \
			--proto_path=api/proto/v1 --proto_path=third_party/ \
			--go_out=plugins=grpc:pkg/api/v1 $$file $$i || exit 1 ;\
	done
	@ echo "Compiling protobufs to go - success \n"

restgateway:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@ echo "Compiling protobufs to restgateway..."
	@ protoc \
 		--proto_path=api/proto/v1 --proto_path=third_party/ \
		  --grpc-gateway_out=logtostderr=true:pkg/api/v1 $(PROTO_ENTRY);
	@ echo "Compiling protobufs to restgateway - success \n"


swagger:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@ echo "Compiling protobufs to swagger definitions..."
	@ protoc \
 		--proto_path=api/proto/v1 --proto_path=third_party/ \
 		--swagger_out=logtostderr=true:api/swagger/v1 $(PROTO_ENTRY)
	@ echo "Compiling protobufs to swagger - success \n"

# Create release binary for current Platform
releaseb:
	@ echo "Compiling release binary"
	@ go build -ldflags "-s -w -X main.gitCommit=$(LAST_COMMIT) -X main.gitTag=$(TAG)" -o dist/server-$(PLATFORM_INFO) cmd/server/main.go;
	@ echo "Compiling binary success - output=dist/server-$(PLATFORM_INFO)"

# Create release binary for current Platform
releaseclientb:
	@ echo "Compiling client release binary"
	@ go build -ldflags "-s -w" -o dist/client-$(PLATFORM_INFO) cmd/client-grpc/main.go;
	@ echo "Compiling client binary success - output=dist/client-$(PLATFORM_INFO)"

# Create release binary for current Platform
binaryfordocker:
	@ echo "Compiling client release binary"
	@ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.gitCommit=$(LAST_COMMIT) -X main.gitTag=$(TAG)" -a -installsuffix cgo -o /go/bin/onionlife ./cmd/server/main.go
	@ echo "Compiling client binary success - output=dist/server-$(PLATFORM_INFO)"

# Create debug binary for current platform
debugb:
	go build -o debug cmd/server/main.go;
