VERSION_FILE := ./internal/versions/versions.go
PROJECT_NAME := oidc-mock
VERSION ?= 0.13
REVISION := $(shell git describe --match=NeVeRmAtCh --always --abbrev=8 --dirty)
PROJECT_SHA := $(shell git rev-parse HEAD)
PROJECT_SHA_SHORT := $(shell git rev-parse HEAD | cut -c1-8)
PROJECT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
DOCKER_REGISTRY ?= gcr.io/aporetodev
DOCKER_IMAGE_NAME ?= $(PROJECT_NAME)
ifeq "$(PROJECT_BRANCH)" "master"
	DOCKER_IMAGE_TAG := latest
	LOCALHOST_SAN :=
else
	DOCKER_IMAGE_TAG := $(REVISION)
	LOCALHOST_SAN := --dns localhost --ip 127.0.0.1
endif
TG := go.aporeto.io/tg
GENCREDS := scripts/gencreds.sh
MKVERSION := scripts/mkversion.sh
SRC := $(shell find . -name .history -prune -o -name \*.go -print | grep -v $(VERSION_FILE) )

GOMODULES := GO111MODULES=on

.PHONY: all version versions build .data package docker docker_build docker_push clean .data

all:
	@ echo "Make targets are 'build', 'version',  'docker_build' or 'docker_push':"
	@ echo
	@ echo "'make build'        - make the binary"
	@ echo "'make versions'     - rebuild the versions pkg files"
	@ echo "'make docker'       - make the binary, and build the docker container (same as 'make docker_build)'"
	@ echo "'make docker_build' - make the binary, and build the docker container (same as 'make docker')"
	@ echo "'make docker_push'  - make the binary and build and push the docker container to GCR tagged as ':$(DOCKER_IMAGE_TAG)'"

version: versions
versions:
	mkdir -p $$(dirname $(VERSION_FILE) )
	$(MKVERSION) "$(VERSION_FILE)"  "$(VERSION)" "$(PROJECT_NAME)" "$(PROJECT_SHA)" "$(PROJECT_BRANCH)" "$(REVISION)"

$(VERSION_FILE): Makefile $(MKVERSION) $(SRC)
	@ $(MAKE) version

build: oidcmock

oidcmock: $(VERSION_FILE) $(SRC) .data
	env $(GOMODULES) go build -o oidcmock

oidcmock.386: $(VERSION_FILE) $(SRC) .data
	env $(GOMODULES) GOOS=linux GOARCH=386 go build -o oidcmock.386

.data:
	go get $(TG)
	$(GENCREDS) $(LOCALHOST_SAN) --dns oidcmock.aporeto.us --dns apotests.oidc.aporeto.us --force

package: .data oidcmock.386
	cp oidcmock.386 docker/oidcmock
	rm -rf docker/.data
	cp -a .data docker/.data

docker: docker_build

docker_build: package
		cd docker \
		  && docker \
		     build --no-cache=true \
			 -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) . \
		  && docker \
		     tag $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) \
			     $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(PROJECT_BRANCH) \
		  && echo Successfully tagged $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(PROJECT_BRANCH)

docker_push: docker_build
		docker \
			push \
			$(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
		docker \
			push \
			$(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(PROJECT_BRANCH)

clean:
	rm -rf vendor
	rm -rf .data
	rm -rf Gopkg.lock
	rm -rf oidcmock oidcmock.386
	rm -rf docker/oidcmock
	rm -rf docker/.data
	rm -rf $(VERSION_FILE)

