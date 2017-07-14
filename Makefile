################
## Important/Override-able bins.
GO ?= go
DOCKER ?= docker
NPM ?= npm
NODE ?= node
RAML2HTML ?= raml2html


################
## Source & Binary config

# This directory.
SOURCEDIR = .

# If a go file changes, Make will see it as dirty under this, and rebuild.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

# A binary folder will always have a main.go, so we're gonna look for those.
# The regex/sed repeats the folder name to denote the binary name Make should look for.
BINARIES := $(shell find . -type f -name 'main.go' ! -path './vendor/*' | sed -n 's/^..\([a-z]*\)\/main.go/.\/\1\/\1/gp')

# The name of this repo in case things need it below.
SL = gitlab.com/packetgg/hearth


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
## JS definitions

# Find all package.jsons
JSTARGETS := $(dir $(shell find $(SOURCEDIR) -name 'package.json' ! -path './vendor/*' ! -path '*/node_modules/*'))

# JS sources
JSSOURCES := $(shell find $(SOURCEDIR) \( -name '*.js' -or -name '*.vue' \) ! -path './vendor/*' ! -path '*/node_modules/*')



################
## RAML definitions

# Find all .ramls
RAMLSOURCES := $(shell find $(SOURCEDIR) -name '*.raml')

# Output the HTML
RAMLTARGETS := $(RAMLSOURCES:%.raml=%.doc.html)

# RAML flags
RAMLFLAGS :=


################
## Docker vars

# Our repo org root
DOCKER_TAG_PREFIX = us.gcr.io/packetgg/hearth/

# The image tag, hopefully just the git hash.
DOCKER_TAG_SUFFIX ?= :$(HASH)$(DIRTY)

# Lists the general images for Docker build targets.
DOCKER_GENERAL := $(shell ls -d misc/dockerfiles/* | sed 's/misc\/dockerfiles\///g')

RETHINKDB_ADDR ?= $(shell docker-compose port rethink 28015)


##############
### TARGETS ###
################
default: all
all: gen js raml bin

################
## Go code

# Builds a Linux binary of a target cmd, then builds it's Docker image.
# The Docker image will be named without the `sl-` prefix, e.g. `sl-launcher` => `launcher`
bin: $(BINARIES)
$(BINARIES): NAME = $(notdir $@)
$(BINARIES): $(SOURCES)
	env GOOS=linux $(GO) build $(LDFLAGS) -o ./$@ -v ./$(dir $@)
	@if [ -e 'Dockerfile' ]; then\
	  $(DOCKER) build ./$(dir $@) -t $(DOCKER_TAG_PREFIX)$(NAME)$(DOCKER_TAG_SUFFIX);\
	fi

.PHONY: test
test:
	env RETHINKDB_ADDR=$(RETHINKDB_ADDR) $(GO) test $(shell glide nv)
	env RETHINKDB_ADDR=$(RETHINKDB_ADDR) overalls -project=$(SL) -covermode=set
	go tool cover -func=overalls.coverprofile

################
## JavaScript
js: $(JSTARGETS)
.PHONY: $(JSTARGETS)
$(JSTARGETS): $(JSSOURCES)
	cd $@ && \
	$(NPM) install && \
	$(NPM) run build

################
## RAML
raml: $(RAMLTARGETS)
$(RAMLTARGETS):
	$(RAML2HTML) $(RAMLFLAGS) $(@:%.doc.html=%.raml) > $@

################
## Codegen
.PHONY: gen
gen:
	go generate $(shell glide nv)

################
## Docker & Kubernetes

# Builds the Docker images that aren't used directly for holding binaries.
.PHONY: $(DOCKER_GENERAL)
$(DOCKER_GENERAL):
	$(DOCKER) build misc/dockerfiles/$@ -t $(DOCKER_TAG_PREFIX)$@$(DOCKER_TAG_SUFFIX)
	$(DOCKER) push $(DOCKER_TAG_PREFIX)$@$(DOCKER_TAG_SUFFIX)

################
## Utilities
clean-all: clean

.PHONY: clean
clean:
	-rm -f $(BINARIES)
	-rm -f $(RAMLTARGETS)