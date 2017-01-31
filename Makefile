GO ?= go
PROTOC ?= protoc

SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

SL = github.com/kayteh/saving-light
HASH := $(shell test -d .git && git rev-parse --short HEAD || echo "UNKNOWN")
DIRTY := $(shell git diff --exit-code >/dev/null || echo "-dirty" && echo "")
BUILD_DATE := $(shell date -u +%FT%T%z)

LDFLAGS = -ldflags "-X ${SL}/etc.Ref=${HASH}${DIRTY} -X ${SL}/etc.BuildDate=${BUILD_DATE}"

BINARIES := $(shell ls -d cmd/* | sed 's/cmd\///g')
PROTOBUF := $(shell find $(SOURCEDIR)/cmd -name '*.proto')
PROTOTARGETS := $(PROTOBUF:.proto=.pb.go)

.PHONY: clean
all: bin grpc
bin: $(BINARIES) 
grpc: $(PROTOTARGETS)

$(BINARIES): grpc
$(BINARIES): $(SOURCES)
	$(GO) build $(LDFLAGS) -v ./cmd/$@/$@.go

$(PROTOTARGETS): $(PROTOBUF)
	$(PROTOC) -I $(dir $@) $(@:.pb.go=.proto) --go_out=plugins=grpc:$(dir $@)

clean:
	-rm -f $(BINARIES) $(PROTOTARGETS)