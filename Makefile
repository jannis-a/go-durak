OUTDIR = _output
COVERAGE = $(OUTDIR)/.$$$$.cov

GOPATH = $(shell go env GOPATH)
GOFLAGS = GOPATH=$(GOPATH) GOBIN=$(OUTDIR)
GOFLAGS_PROD = GOPATH=$(GOPATH) CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GOCMD = $(GOFLAGS) go

DOCKER_IMAGE = jannis-a/go-durak

#
# Helper targets
#
output-dir:
	mkdir -p $(OUTDIR)

build: output-dir
	$(info Building binaries...)
	$(GOCMD) list ./cmd/... | { cd $(OUTDIR) && $(GOFLAGS) xargs -n 1 -- go build -v; }

build-prod: output-dir # TODO: test
	$(info Building release binaries...)
	$(GOCMD) list ./cmd/... | { cd $(OUTDIR) && $(GOFLAGS_PROD) xargs -n 1 -- go build -v; }

docker:
	$(info Building docker container...)
	docker build . -t $(DOCKER_IMAGE)

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

clean:
	$(info Cleaning files...)
	rm -rf $(OUTDIR) vendor
