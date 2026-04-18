BINARY   := mand
CMD      := ./cmd/main.go
PREFIX   ?= /usr/local

.PHONY: build install uninstall test fmt lint clean

build:
	go build -o $(BINARY) $(CMD)

install: build
	install -d $(PREFIX)/bin
	install -m 755 $(BINARY) $(PREFIX)/bin/$(BINARY)

uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, falling back to go vet"; \
		go vet ./...; \
	fi

clean:
	rm -f $(BINARY)
