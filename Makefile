DOCKER ?= docker

CROSSBUILD_OS   = linux windows darwin
CROSSBUILD_ARCH = 386 amd64

VERSION   := $(shell git describe --tags --dirty=-dirty)
REVISION  := $(shell git describe --abbrev=0 --always --match=always-commit-hash --dirty=-dirty)
BRANCH    := $(shell git rev-parse --abbrev-ref HEAD)
BUILDDATE := $(shell date --iso-8601=seconds)
BUILDUSER ?= $(USER)
BUILDHOST ?= $(HOSTNAME)
LDFLAGS    = -X github.com/prometheus/common/version.Version=$(VERSION) \
             -X github.com/prometheus/common/version.Revision=$(REVISION) \
             -X github.com/prometheus/common/version.Branch=$(BRANCH) \
             -X github.com/prometheus/common/version.BuildUser=$(BUILDUSER)@$(BUILDHOST) \
             -X github.com/prometheus/common/version.BuildDate=$(BUILDDATE)

all: build test

test:
	@echo ">> testing code"
	@go test -v -cover ./...

build:
	@echo ">> building binaries"
	@go build -ldflags="$(LDFLAGS)"

crossbuild: $(GOPATH)/bin/gox
	@echo ">> cross-building"
	@gox -arch="$(CROSSBUILD_ARCH)" -os="$(CROSSBUILD_OS)" -ldflags="$(LDFLAGS)" -output="binaries/stream_exporter_{{.OS}}_{{.Arch}}"

release: bin/github-release
	@echo ">> uploading release ${VERSION}"
	@for bin in binaries/*; do \
		./bin/github-release upload -t ${VERSION} -n $$(basename $${bin}) -f $${bin}; \
	done

docker:
	@echo ">> building docker image"
	@$(DOCKER) build -t carlpett/stream_exporter .

$(GOPATH)/bin/gox:
	# Need to disable modules for this to not pollute go.mod
	@GO111MODULE=off go get -u github.com/mitchellh/gox

bin/github-release:
	@mkdir -p bin
	@curl -sL 'https://github.com/aktau/github-release/releases/download/v0.6.2/linux-amd64-github-release.tar.bz2' | tar xjf - --strip-components 3 -C bin

.PHONY: all build crossbuild test release
