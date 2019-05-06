VERSION_FILE := ./version/version.go
PROJECT_NAME := oidc-mock
BUILD_NUMBER := latest
VERSION := 0.11
REVISION=$(shell git log -1 --pretty=format:"%H")
DOCKER_REGISTRY?= gcr.io/aporetodev
DOCKER_IMAGE_NAME?=$(PROJECT_NAME)
DOCKER_IMAGE_TAG?=$(BUILD_NUMBER)

build:
	env GOOS=linux GOARCH=386 go build -o oidcmock

package: build
	mv oidcmock docker/oidcmock
	rm -rf docker/.data
	cp -a .data docker/.data

clean:
	rm -rf vendor
	rm -rf .data
	rm -rf Gopkg.lock

docker_build: package
		cd docker && docker \
			build --no-cache=true --rm \
			-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

docker_push: docker_build
		docker \
			push \
			$(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
