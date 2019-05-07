COVERAGE = coverage.out
API = durak

build:
	go build $(API).go

test:
	go test ./...

test-coverage:
	go test -coverprofile=$(COVERAGE) ./...

coverage: test-coverage
	go tool cover -func=$(COVERAGE)

coverage-html: test-coverage
	go tool cover -html=$(COVERAGE)

clean:
	rm $(API) $(COVERAGE)
	
#build-linux: test
#    $(info Building Linux Binaries...)
#    mkdir -p $(OUTDIR)
#    $(GOCMD) list ./cmd/... | { cd $(OUTDIR) && $(GOFLAGS_LINUX_AMD64) xargs -n 1 -- go build -v; }
#
#GOPATH = $(shell go env GOPATH)
#GOFLAGS = GOPATH=$(GOPATH) GOBIN=$(OUTDIR)
#GOFLAGS_LINUX_AMD64 = GOPATH=$(GOPATH) CGO_ENABLED=0 GOOS=linux GOARCH=amd64
#GOCMD = $(GOFLAGS) go