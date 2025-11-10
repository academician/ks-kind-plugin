DATE=$(shell date -u +%Y-%m-%d)
VERSION=$(shell cat VERSION | sed 's/-dev//g')

.PHONY: format
format: $(GOLICENSES) $(GOIMPORTS)
	goimports -l -w .

.PHONY: check
check:
	golangci-lint run "${SOURCE_TREES[.]}" --timeout=10m0s --verbose --print-resources-usage --modules-download-mode=vendor

.PHONY: build
build:
	@go build -ldflags "-w -X github.com/academician/ks-kind-plugin.version=${VERSION} -X github.com/academician/ks-kind-plugin.buildDate=${DATE}" -o ks-kind-plugin .

.PHONY: all
all: format check build

.PHONY: revendor
revendor:
	@GO111MODULE=on go mod vendor
	@GO111MODULE=on go mod tidy
