GO ?= go
PROTOC ?= protoc
DOCKER ?= docker

SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

SL = github.com/kayteh/saving-light
HASH := $(shell test -d .git && git rev-parse --short HEAD || echo "UNKNOWN")
DIRTY := $(shell git diff --exit-code >/dev/null || echo "-dirty" && echo "")
BUILD_DATE := $(shell date -u +%FT%T%z)

LDFLAGS = -ldflags "-X ${SL}/etc.Ref=${HASH}${DIRTY} -X ${SL}/etc.BuildDate=${BUILD_DATE}"

BINARIES := $(shell ls -d cmd/* | sed 's/cmd\(.*\)/cmd\1\1/g')
PROTOBUF := $(shell find $(SOURCEDIR)/cmd -name '*.proto')
PROTOTARGETS := $(PROTOBUF:.proto=.pb.go)

DOCKER_TAG_PREFIX = quay.io/saving-light/
DOCKER_TAG_SUFFIX ?= :$(HASH)$(DIRTY)

# General Dockerfiles
DOCKER_GEN := $(shell ls -d misc/dockerfiles/* | sed 's/misc\/dockerfiles\///g')

PROTOBUF_INC := /usr/local/include . ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis ./vendor

.PHONY: clean clean-rpc $(DOCKER_GEN)
default: all
all: grpc bin
clean-all: clean clean-rpc
bin: $(BINARIES) 
grpc: $(PROTOTARGETS)

$(BINARIES): NAME = $(notdir $@)
$(BINARIES): $(SOURCES)
	env GOOS=linux $(GO) build $(LDFLAGS) -o ./$@ -v ./$(dir $@)
	$(DOCKER) build ./$(dir $@) -t $(DOCKER_TAG_PREFIX)$(NAME:sl-%=%)$(DOCKER_TAG_SUFFIX)

$(PROTOTARGETS): $(PROTOBUF)
	$(PROTOC) \
		$(addprefix -I, $(PROTOBUF_INC)) \
		-I $(dir $@) \
		--grpc-gateway_out=logtostderr=true:. \
		--swagger_out=logtostderr=true:. \
		--go_out=M'google/api/annotations.proto'="${SL}/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api",plugins=grpc:. \
		$(@:.pb.go=.proto)

test-docker: 
	@echo $(BINARIES)

$(DOCKER_GEN):
	$(DOCKER) build misc/dockerfiles/$@ -t $(addprefix $(DOCKER_TAG_PREFIX), $(addsuffix $(DOCKER_TAG_SUFFIX), $@))
	$(DOCKER) push $(addprefix $(DOCKER_TAG_PREFIX), $(addsuffix $(DOCKER_TAG_SUFFIX), $@))

clean:
	-rm -f $(BINARIES)

clean-rpc:
	-rm -f $(PROTOTARGETS)