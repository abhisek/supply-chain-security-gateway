all: clean setup server

setup:
	mkdir -p out

server:
	go build -o out/pdp-server cmd/pdp/pdp.go
	go build -o out/tap-server cmd/tap/tap.go
	go build -o out/dcs-server cmd/dcs/dcs.go

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	-rm -rf out
