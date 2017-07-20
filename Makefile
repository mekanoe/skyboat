################
## Important/Override-able bins.
GO ?= go
DOCKER ?= docker
COMPOSE ?= docker-compose -H /var/run/docker.sock
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
SL = skyboat.io/x


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
DOCKER_TAG_PREFIX = skyboat/

# The image tag, hopefully just the git hash.
DOCKER_TAG_SUFFIX ?= :$(HASH)$(DIRTY)

DOCKER_PUSH ?= 0

# Lists the general images for Docker build targets.
DOCKER_GENERAL := $(shell ls -d misc/dockerfiles/* | sed 's/misc\/dockerfiles\///g')

RETHINKDB_ADDR ?= $(shell docker-compose port rethink 28015)


##############
### TARGETS ###
################
default: all
all: js raml bin

################
## Go code

# Builds a Linux binary of a target cmd, then builds it's Docker image.

bin: $(BINARIES)
$(BINARIES): NAME = $(notdir $@)
$(BINARIES): $(SOURCES)
	@echo "🛠 building $@ codegen"	
	@$(GO) generate ./$(dir $@)...
	@if [ `docker-compose ps | grep -c build-server` == "0" ]; then\
	  echo "> starting build server";\
	  docker-compose up -d build-server;\
	fi
	@echo "🛠 building $@ binary"
	@docker-compose exec build-server ash -c '\
	  cd /go/src/$(SL); \
	  $(GO) install $(LDFLAGS) -v ./$(dir $@) &&\
	  mv /go/bin/$(NAME) ./$(NAME)'
	@if [ -e ./$(dir $@)/Dockerfile ]; then\
	  echo "🛠 building $@ docker image $(DOCKER_TAG_PREFIX)$(NAME)$(DOCKER_TAG_SUFFIX)";\
	  $(DOCKER) build ./$(dir $@) -q -t $(DOCKER_TAG_PREFIX)$(NAME)$(DOCKER_TAG_SUFFIX);\
	  if [ $(DOCKER_PUSH) == "1" ]; then\
	   echo "📤 pushing $(DOCKER_TAG_PREFIX)$(NAME)$(DOCKER_TAG_SUFFIX)";\
	   $(DOCKER) push $(DOCKER_TAG_PREFIX)$(NAME)$(DOCKER_TAG_SUFFIX) >/dev/null;\
	  else\
	   echo "💤 not pushing $(DOCKER_TAG_PREFIX)$(NAME)$(DOCKER_TAG_SUFFIX)";\
	  fi;\
	fi
	@echo "✅ done!"

.PHONY: test
test:
	$(GO) test $(shell glide nv)
	overalls -project=$(SL) -covermode=set
	go tool cover -func=overalls.coverprofile

.PHONY: coverage
coverage: test
coverage:
	go tool cover -html=overalls.coverprofile
	

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
	go generate `glide nv`

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


.PHONY: tmpl
tmpl:
	bash ./tools/tmpl.bash $(filter-out $@,$(MAKECMDGOALS))

.PHONY: kup
kup:
	@if [ `helm list | grep -c skyboat-dev` == "0" ]; then\
	  helm install ./charts --namespace skyboat --name skyboat-dev --set docker_tag=$(HASH)$(DIRTY),repo_prefix=$(DOCKER_TAG_PREFIX);\
	else\
	  helm upgrade skyboat-dev ./charts --namespace skyboat --set docker_tag=$(HASH)$(DIRTY),repo_prefix=$(DOCKER_TAG_PREFIX);\
	fi

.PHONY: reset-env
reset-env:
	docker-compose down
	-rm -rf .cache
	helm delete --purge skyboat-dev

restotest: $(shell find ./restokit -name "*.go")
	$(GO) generate ./restokit/...
	$(GO) run ./restokit/restotest/main.go
