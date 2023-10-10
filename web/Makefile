.PHONY: build
build:
	go build -o app

.PHONY: build_prod
build_prod: generate_prod
	GOOS=linux GOARCH=amd64 go build -o app .

.PHONY: unit_tests
unit_tests: generate
	gotestsum --hide-summary=skipped -- -v ./...

.PHONY: integration_tests
integration_tests: generate_integration_test
	gotestsum --hide-summary=skipped -- -tags integration ./... -v

.PHONY: tests
tests: unit_tests

.PHONY: coverage
coverage:
	go test -tags integration -coverpkg=./... -coverprofile=coverage.o ./... > /dev/null
	go tool cover -html=coverage.o -o coverage.html

.PHONY: coverage_func
coverage_func:
	go test -tags integration -coverpkg=./... -coverprofile=coverage.o ./... > /dev/null
	go tool cover -func coverage.o

.PHONY: generate
generate: clean
	go generate -x -tags wireinject ./...

.PHONY: generate_prod
generate_prod:
	go generate -x -tags wireinject wire/wire.go

.PHONY: generate_integration_test
generate_integration_test: clean
	go generate -x -tags="wireinject integration" ./...

.PHONY: clean
clean:
	git ls-files --others --ignored --exclude-standard ./ | \
	grep -v '.env*' | \
	grep -v '.idea*' | \
	xargs -I {} rm -f {}