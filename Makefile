COVERAGE = coverage.out
API = durak

build:
	go build $(API).go

test:
	go test ./...

test-coverage:
	go test -coverprofile=$(COVERAGE) ./...

coverage:
	go tool cover -func=$(COVERAGE)

coverage-html:
	go tool cover -html=$(COVERAGE)

clean:
	rm $(API) $(COVERAGE)