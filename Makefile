.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o bin/regex-tui cmd/regex-tui/main.go

.PHONY: clean
clean:
	rm -f bin/regex-tui

.PHONY: install
install:
	cp bin/regex-tui /usr/local/bin/regex-tui

.PHONY: uninstall
uninstall:
	rm -f /usr/local/bin/regex-tui

.PHONY: lint
lint:
	go vet ./...
	go fmt ./...

.PHONY: modernize
modernize:
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

.PHONY: run
run:
	go run cmd/regex-tui/main.go

.PHONY: demo
demo:
	cd assets && vhs demo.tape