name = imis

all: clean build

clean:
	- go clean
	- go get

build:
	- go build
	- docker build . --tag $(name)

run:
	docker run --name $(name) -p 3000:3000/tcp --restart unless-stopped -dti $(name)

.PHONY: clean build