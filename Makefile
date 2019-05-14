OUTDIR = _output
CONFIG = config.yml
COVERAGE = $(OUTDIR)/.$$$$.cov

GOPATH = $(shell go env GOPATH)
GOFLAGS = GOPATH=$(GOPATH) GOBIN=$(OUTDIR) GO111MODULE=on
GOFLAGS_PROD = GOPATH=$(GOPATH) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on
GOCMD = $(GOFLAGS) go

DOCKER_IMAGE = jannis-a/go-durak

#
# Helpers
#
output-dir:
	mkdir -p $(OUTDIR)

copy-config:
	cp config-example.yml $(CONFIG)

#
# Build
#
build: output-dir
	$(info Building binaries...)
	$(GOCMD) list ./cmd/... | { cd $(OUTDIR) && $(GOFLAGS) xargs -n 1 -- go build -v; }

build-prod: output-dir # TODO: test
	$(info Building production binaries...)
	$(GOCMD) list ./cmd/... | { cd $(OUTDIR) && $(GOFLAGS_PROD) xargs -n 1 -- go build -v; }

docker:
	$(info Building docker container...)
	docker build . -t $(DOCKER_IMAGE)

#
# Test and coverage
#
test:
	$(info Running tests...)
	$(GOCMD) test ./...

test-coverage: output-dir
	$(info Running tests with coverage...)
	$(GOCMD) test -coverprofile=$(COVERAGE) ./...

coverage: test-coverage
	$(GOCMD) tool cover -func=$(COVERAGE)

coverage-html: test-coverage
	$(GOCMD) tool cover -html=$(COVERAGE)

#
# Clean
#
clean:
	$(info Cleaning build files...)
	rm -rf $(OUTDIR)

clean-all: clean
	rm $(CONFIG)
