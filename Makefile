################
# Important/Override-able bins.
GO ?= go
PROTOC ?= protoc
DOCKER ?= docker

################
# Source & Binary config
SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
BINARIES := $(shell ls -d cmd/* | sed 's/cmd\(.*\)/cmd\1\1/g')
SL = github.com/kayteh/saving-light

################
# Build Date & Git Revision
HASH := $(shell test -d .git && git rev-parse --short HEAD || echo "UNKNOWN")
DIRTY := $(shell git diff --exit-code >/dev/null || echo "-dirty" && echo "")
BUILD_DATE := $(shell date -u +%FT%T%z)
LDFLAGS = -ldflags "-X ${SL}/etc.Ref=${HASH}${DIRTY} -X ${SL}/etc.BuildDate=${BUILD_DATE}"


################
# gRPC/Protobuf definitions
PROTOBUF := $(shell find $(SOURCEDIR)/cmd -name '*.proto')
PROTOBUF_INC := /usr/local/include . ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis ./vendor
PROTOTARGETS := $(PROTOBUF:.proto=.pb.go)

################
# Docker vars
DOCKER_TAG_PREFIX = quay.io/saving-light/
DOCKER_TAG_SUFFIX ?= :$(HASH)$(DIRTY)
DOCKER_GENERAL := $(shell ls -d misc/dockerfiles/* | sed 's/misc\/dockerfiles\///g')


##############
### TARGETS ###
################
default: all
all: grpc bin

################
# Go code
bin: $(BINARIES) 
$(BINARIES): NAME = $(notdir $@)
$(BINARIES): $(SOURCES)
	env GOOS=linux $(GO) build $(LDFLAGS) -o ./$@ -v ./$(dir $@)
	$(DOCKER) build ./$(dir $@) -t $(DOCKER_TAG_PREFIX)$(NAME:sl-%=%)$(DOCKER_TAG_SUFFIX)

################
# gRPC/Protobuf 
grpc: $(PROTOTARGETS)
$(PROTOTARGETS): $(PROTOBUF)
	$(PROTOC) \
		$(addprefix -I, $(PROTOBUF_INC)) \
		-I $(dir $@) \
		--grpc-gateway_out=logtostderr=true:. \
		--swagger_out=logtostderr=true:. \
		--go_out=M'google/api/annotations.proto'="${SL}/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api",plugins=grpc:. \
		$(@:.pb.go=.proto)

################
# Docker & Kubernetes
.PHONY: $(DOCKER_GENERAL)
$(DOCKER_GENERAL):
	$(DOCKER) build misc/dockerfiles/$@ -t $(DOCKER_TAG_PREFIX)$@$(DOCKER_TAG_SUFFIX)
	$(DOCKER) push $(DOCKER_TAG_PREFIX)$@$(DOCKER_TAG_SUFFIX)

################
# Utilities
clean-all: clean clean-rpc

.PHONY: clean
clean:
	-rm -f $(BINARIES)

.PHONY: clean-rpc
clean-rpc:
	-rm -f $(PROTOTARGETS)