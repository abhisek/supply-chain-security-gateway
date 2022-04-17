all: clean setup server

setup:
	mkdir -p out

server:
	go build -o out/pdp-server pdp.go

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	-rm -rf out

