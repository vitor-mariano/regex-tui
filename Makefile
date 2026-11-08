.DEFAULT_GOAL := bin/regex-tui

GOPATH := $(shell go env GOPATH)
VERSION ?= master

PLATFORMS := linux darwin windows
ARCHITECTURES := amd64 arm64

bin/regex-tui:
	go build -o bin/regex-tui main.go

.PHONY: clean
clean:
	rm -f bin/*

.PHONY: release
release: clean
	@for platform in $(PLATFORMS); do \
		for arch in $(ARCHITECTURES); do \
			ext=""; \
			if [ "$$platform" = "windows" ]; then ext=".exe"; fi; \
			output="bin/regex-tui_$(VERSION)_$${platform}.$${arch}$${ext}"; \
			echo "Building $$output..."; \
			GOOS=$$platform GOARCH=$$arch go build -o $$output main.go; \
		done; \
	done

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

.PHONY: test
test:
	go test ./...

# Refresh testdata/*.golden snapshots after intentional UI changes.
.PHONY: test-update
test-update:
	go test ./pkg/components/regexview/ -update

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