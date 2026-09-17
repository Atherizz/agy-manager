ifeq ($(OS),Windows_NT)
	BINARY := agym.exe
	BIN_DIR := $(USERPROFILE)/bin
	RM := del /Q
	DEVNULL := nul
	CP := copy /Y
	MKDIR := mkdir
else
	BINARY := agym
	BIN_DIR := $(HOME)/bin
	RM := rm -f
	DEVNULL := /dev/null
	CP := cp -f
	MKDIR := mkdir -p
endif

VERSION := $(shell git describe --tags --always --dirty 2>$(DEVNULL) || echo "0.1.0")
COMMIT  := $(shell git rev-parse --short HEAD 2>$(DEVNULL) || echo "unknown")
LDFLAGS := -ldflags "-s -w -X github.com/Atherizz/agy-manager/cmd.version=$(VERSION) -X github.com/Atherizz/agy-manager/cmd.commit=$(COMMIT)"

.PHONY: build install clean test

build:
	go build $(LDFLAGS) -o $(BINARY) .

install: build
	@$(MKDIR) $(BIN_DIR) 2>$(DEVNULL) || true
	$(CP) $(BINARY) $(BIN_DIR)/$(BINARY)

clean:
	$(RM) $(BINARY) 2>$(DEVNULL) || true

test:
	go test -v ./...
