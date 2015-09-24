default: build

prepare:
	go get ./...

test: prepare
	go test ./...

build: test
	go build ./cmd/packer-provisioner-serverspec

install: build
	mkdir -p ~/.packer.d/plugins
	install ./packer-provisioner-serverspec ~/.packer.d/plugins/

.PHONY: default prepare test build install
