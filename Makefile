VERSION ?= 0.1.0

COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -ldflags "-X github.com/Atherizz/agy-manager/cmd.version=$(VERSION) -X github.com/Atherizz	/agy-manager/cmd.commit=$(COMMIT)"

  

build:

    go build $(LDFLAGS) -o agym.exe .
