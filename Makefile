VERSION := $(shell git describe --tags --always --dirty 2>nul || echo "0.1.0")
COMMIT  := $(shell git rev-parse --short HEAD 2>nul || echo "unknown")
LDFLAGS := -ldflags "-s -w -X github.com/Atherizz/agy-manager/cmd.version=$(VERSION) -X github.com/Atherizz/agy-manager/cmd.commit=$(COMMIT)"

.PHONY: build install clean

build:
	go build $(LDFLAGS) -o agym.exe .

install: build
	copy /Y agym.exe "%USERPROFILE%\bin\agym.exe"

clean:
	del /Q agym.exe 2>nul || true
