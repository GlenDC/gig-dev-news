OUTPUT ?= bin
GOOS ?= linux
GOARCH ?= amd64

PACKAGE = github.com/glendc/gig-dev-news
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE = $(shell date +%FT%T%z)

ldflags = -extldflags "-static"
ldflagsversion = -X main.CommitHash=$(COMMIT_HASH) -X main.BuildDate=$(BUILD_DATE) -s -w

all: gdnbot

gdnbot: $(OUTPUT)
ifeq ($(GOOS), darwin)
	GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags '$(ldflagsversion)' -o $(OUTPUT)/gdnbot $(PACKAGE)
else
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags '$(ldflags)$(ldflagsversion)' -o $(OUTPUT)/gdnbot $(PACKAGE)
endif

$(OUTPUT):
	mkdir -p $(OUTPUT)

.PHONY: gdnbot $(OUTPUT)
