IMAGE := json0/terraform-upgrade
VERSION := 0.0.5
TAGGED_IMAGE := $(IMAGE):$(VERSION)

build: .built-docker-image

.built-docker-image: Dockerfile makefile
	docker build -t $(IMAGE) .
	touch .built-docker-image

tag: .built-docker-image
	docker tag $(IMAGE) $(TAGGED_IMAGE)
	docker tag $(IMAGE) $(IMAGE):latest

push:
	make tag
	docker push $(TAGGED_IMAGE)
	docker push $(IMAGE):latest

all: 
	make push

shell:
			docker pull $(TAGGED_IMAGE)
			docker run --rm -it \
	-e GITHUB_AUTH_TOKEN=$${GITHUB_AUTH_TOKEN} \
	-e GITHUB_AUTH_USER=$${GITHUB_AUTH_USER} \
	-e SSH_AUTH_SOCK=$SSH_AUTH_SOCK \
							-v $$(pwd):/app \
							-v $${HOME}/.ssh:$${HOME}/.ssh \
							-v $${HOME}/.gitconfig:$${HOME}/.gitconfig \
							-v $${HOME}/.config/gh:$${HOME}/.config/gh \
							-v $(dirname $SSH_AUTH_SOCK):$(dirname $SSH_AUTH_SOCK) \
							-w /app \
							$(TAGGED_IMAGE) bash

run:
			docker pull $(TAGGED_IMAGE)
			docker run --rm -it \
	-e GITHUB_AUTH_TOKEN=$${GITHUB_AUTH_TOKEN} \
	-e GITHUB_AUTH_USERNAME=$${GITHUB_AUTH_USERNAME} \
							-v $$(pwd):/app \
							-v $${HOME}/.ssh:/root/.ssh \
							-w /app \
							$(TAGGED_IMAGE) go run main.go

clean:
	rm -rf cloud-platform-*
