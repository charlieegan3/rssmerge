PROJECT := rssmerge
TAG := $(shell tar -cf - . | md5sum | cut -f 1 -d " ")

build:
	docker build -t charlieegan3/$(PROJECT):${TAG} .
	docker build -t charlieegan3/$(PROJECT):arm-${TAG} -f Dockerfile.arm .

push: build
	docker push charlieegan3/$(PROJECT):${TAG}
	docker push charlieegan3/$(PROJECT):arm-${TAG}
