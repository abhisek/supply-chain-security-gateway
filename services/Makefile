all: clean setup server

oapi-codegen-install:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

protoc-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest


oapi-codegen:
	oapi-codegen -package openssf -generate types ./spec/openssf/osv-api-openapi.yml > ./pkg/common/openssf/osv.types.go

protoc-codegen:
	protoc -I ./spec/proto/ \
		-I ./spec/proto/lib/ -I ./spec/proto/lib/protoc-gen-validate \
		--go_out=./gen/ \
		--go-grpc_out=./gen/ \
		--validate_out="lang=go:./gen/" \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--validate_opt=paths=source_relative \
		./spec/proto/pds.proto

setup:
	mkdir -p out

server: oapi-codegen
	go build -o out/pdp-server cmd/pdp/pdp.go
	go build -o out/tap-server cmd/tap/tap.go
	go build -o out/dcs-server cmd/dcs/dcs.go

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	-rm -rf out
