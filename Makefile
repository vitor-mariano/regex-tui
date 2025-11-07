.DEFAULT_GOAL := build

GOPATH := $(shell go env GOPATH)

.PHONY: build
build:
	go build -o bin/regex-tui main.go

.PHONY: clean
clean:
	rm -f bin/regex-tui

.PHONY: debug
debug:
	go build -gcflags="-N -l" -o bin/regex-tui main.go
	./bin/regex-tui

.PHONY: install
install:
	go install .

.PHONY: uninstall
uninstall:
	rm -f $(GOPATH)/bin/regex-tui

.PHONY: lint
lint:
	go vet ./...
	go fmt ./...

.PHONY: modernize
modernize:
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

.PHONY: run
run:
	go run main.go

.PHONY: demo
demo:
	cd assets && vhs demo.tape