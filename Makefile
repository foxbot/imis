name = imis

all: clean build

clean:
        - docker stop $(name) && docker rm $(name)
        - go clean
#       - go get

build:
        - CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(name) .
        - docker build . --tag $(name)

run:
        docker run --name $(name) -p 80:3000/tcp --restart unless-stopped -dti $(name)

.PHONY: clean build