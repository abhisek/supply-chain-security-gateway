all: clean setup server

oapi-codegen-install:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

oapi-codegen:
	oapi-codegen -package openssf -generate types ./spec/openssf/osv-api-openapi.yml > ./pkg/common/openssf/osv.types.go

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
