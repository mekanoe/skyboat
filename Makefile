################
## Important/Override-able bins.
GO ?= go
PROTOC ?= protoc
DOCKER ?= docker


################
## Source & Binary config

# This directory.
SOURCEDIR = .

# If a go file changes, Make will see it as dirty under this, and rebuild.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

# We only have compilable binaries inside the /cmd/ directory, so we'll target the top level ones.
# The regex/sed repeats the folder name to denote the binary name Make should look for.
BINARIES := $(shell ls -d cmd/* | sed 's/cmd\(.*\)/cmd\1\1/g')

# The name of this repo in case things need it below.
SL = github.com/kayteh/spaceplane


################
## Build Date & Git Revision

# If this is a git repo (and not some sort of odd archive,) get the rev-hash.
HASH := $(shell test -d .git && git rev-parse --short HEAD || echo "UNKNOWN")

# If code has changed since last commit, mark as dirty.
DIRTY := $(shell git diff --exit-code >/dev/null || echo "-dirty" && echo "")

# Generate an ISO-8601 date.
BUILD_DATE := $(shell date -u +%FT%T%z)

# -ldflags for go build. Overrides some variables so they can be used as version markers.
LDFLAGS = -ldflags "-X ${SL}/etc.Ref=${HASH}${DIRTY} -X ${SL}/etc.BuildDate=${BUILD_DATE}"


################
## gRPC/Protobuf definitions

# Any .proto under /cmd/ is useful to us, this targets them.
PROTOBUF := $(shell find $(SOURCEDIR)/cmd -name '*.proto')

# This is a list of includes we need for gRPC, mostly includes stuff needed for gRPC-gateway.
PROTOBUF_INC := /usr/local/include . ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis ./vendor

# Translates `service.proto` => `service.pb.go` because these are the main artifacts of the gRPC targets.
PROTOTARGETS := $(PROTOBUF:.proto=.pb.go)


################
## Docker vars

# Our repo org root
DOCKER_TAG_PREFIX = quay.io/spaceplane/

# The image tag, hopefully just the git hash.
DOCKER_TAG_SUFFIX ?= :$(HASH)$(DIRTY)

# Lists the general images for Docker build targets.
DOCKER_GENERAL := $(shell ls -d misc/dockerfiles/* | sed 's/misc\/dockerfiles\///g')



##############
### TARGETS ###
################
default: all
all: grpc bin

################
## Go code

# Builds a Linux binary of a target cmd, then builds it's Docker image.
# The Docker image will be named without the `sl-` prefix, e.g. `sl-launcher` => `launcher`
bin: $(BINARIES) 
$(BINARIES): NAME = $(notdir $@)
$(BINARIES): $(SOURCES)
	env GOOS=linux $(GO) build $(LDFLAGS) -o ./$@ -v ./$(dir $@)
	$(DOCKER) build ./$(dir $@) -t $(DOCKER_TAG_PREFIX)$(NAME:spaceplane-%=%)$(DOCKER_TAG_SUFFIX)

################
## gRPC/Protobuf 

# Generates gRPC interfaces and other odds and ends.
# Currently generates: Go, gRPC Gateway (Go), Swagger
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
## Docker & Kubernetes

# Builds the Docker images that aren't used directly for holding binaries.
.PHONY: $(DOCKER_GENERAL)
$(DOCKER_GENERAL):
	$(DOCKER) build misc/dockerfiles/$@ -t $(DOCKER_TAG_PREFIX)$@$(DOCKER_TAG_SUFFIX)
	$(DOCKER) push $(DOCKER_TAG_PREFIX)$@$(DOCKER_TAG_SUFFIX)

################
## Utilities
clean-all: clean clean-rpc

.PHONY: clean
clean:
	-rm -f $(BINARIES)

.PHONY: clean-rpc
clean-rpc:
	-rm -f $(PROTOTARGETS)