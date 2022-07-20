SHELL := /bin/bash # Use bash syntax

test:
	echo "Today is $$DATE"

build:	vendor tls
build:
	BUILD_DATE=$$(date -u '+%Y-%m-%d %H:%M:%S UTC') BUILD_VERSION=$(VERSION) docker-compose build

help:   	## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

build:		## Build the application stack

start:		## Start the application stack
	docker-compose up -d

tail:		## Tail the logs
	docker-compose logs -f

demo:		## Run some curl operations on the running application stack
	scripts/demo.sh

stop:		## Stop the application stack
	docker-compose down

clean:		## Clean up any artifacts created
	rm -rf redis/tls
	docker-compose down --rmi all
	rm -rf vendor

tls:		## Generate TLS certificates
	if [ ! -d redis/tls ] ; then cd redis && ./generate-tls-certs.sh ; fi

build-fast:		## Build the go application binary quickly, without injecting build date or version
	go build .

build-binary:	## Build the go application binary
	go build -a -ldflags "-X 'main.buildDate=$$(date -u '+%Y-%m-%d %H:%M:%S UTC')' -X main.buildVersion=$(VERSION)" .

build-docker:	## Build the docker container
	docker build --build-arg BUILD_DATE="$$(date -u '+%Y-%m-%d %H:%M:%S UTC')" --build-arg BUILD_VERSION=$(VERSION) . -t api

vendor:		## Create the go vendor/ directory
	go mod vendor

.PHONY: clean    
