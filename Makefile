VERSION_FILE := ./internal/versions/versions.go
PROJECT_NAME := oidc-mock
VERSION ?= 0.12
REVISION := $(shell git describe --match=NeVeRmAtCh --always --abbrev=8 --dirty)
PROJECT_SHA := $(shell git rev-parse HEAD)
PROJECT_SHA_SHORT := $(shell git rev-parse HEAD | cut -c1-8)
PROJECT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
DOCKER_REGISTRY ?= gcr.io/aporetodev
DOCKER_IMAGE_NAME ?= $(PROJECT_NAME)
ifeq "$(PROJECT_BRANCH)" "master"
	DOCKER_IMAGE_TAG := latest
else
	DOCKER_IMAGE_TAG := $(PROJECT_SHA_SHORT)
endif
DOCKER_IMAGE_TAG ?= $()
TG := go.aporeto.io/tg
GENCREDS := scripts/gencreds.sh
MKVERSION := scripts/mkversion.sh
SRC := $(shell find . -name .history -prune -o -name \*.go -print | grep -v $(VERSION_FILE) )

.PHONY: all version build build.386 .data package  docker_build docker_push clean .data

all:
	@ echo "Make targets are 'build', 'version',  'docker_build' or 'docker_push':"
	@ echo
	@ echo "'make build'        - make the binary"
	@ echo "'make version'      - rebuild the version file"
	@ echo "'make docker_build' - make the binary, and build the docker container"
	@ echo "'make docker_push'  - make the binary and build and push the docker container to GCR tagged as :$(DOCKER_IMAGE_TAG)"

version:
	mkdir -p $$(dirname $(VERSION_FILE) )
	$(MKVERSION) "$(VERSION_FILE)"  "$(VERSION)" "$(PROJECT_NAME)" "$(PROJECT_SHA)" "$(PROJECT_BRANCH)" "$(REVISION)"

$(VERSION_FILE): Makefile $(MKVERSION) $(SRC)
	@ $(MAKE) version

build: $(VERSION_FILE) $(SRC)
	go build -o oidcmock

build.386: $(VERSION_FILE) $(SRC)
	env GOOS=linux GOARCH=386 go build -o oidcmock.386

.data:
	go get $(TG)
	$(GENCREDS) --dns oidcmock.aporeto.us --force

package: .data build.386
	cp oidcmock.386 docker/oidcmock
	rm -rf docker/.data
	cp -a .data docker/.data

docker_build: package
		cd docker \
		  && docker \
		     build \
			 -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) . \
		  && docker \
		     tag $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) \
			     $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(PROJECT_BRANCH)

docker_push: docker_build
		echo docker \
			push \
			$(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
		echo docker
			push \
			$(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(PROJECT_BRANCH)

clean:
	rm -rf vendor
	rm -rf .data
	rm -rf Gopkg.lock
