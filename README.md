# Go API Template

## Overview

This is my template for a well formed REST API in Go.

This template features:

* [gorilla/mux](https://github.com/gorilla/mux) for routing
* [go.uber.org/zap](https://pkg.go.dev/go.uber.org/zap) for logging (this is new to me, I also want to play with zerolog and some others)
* [spf13/viper](https://github.com/spf13/viper) for configuration
  * Server configuration can be done via file, command line, or environment variables (precidence in that order)
* Basic authentication with [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) for password hashing
* [casbin/casbin](https://github.com/casbin/casbin) for authorization
  * Examples [here](https://github.com/casbin/casbin/tree/master/examples)
* Functions and code blocks are timed and logged
* Testing of the handlers

The methodologies and patterns coded herein are not the only way to craft a well-formed REST API, and there is plenty of room to discuss whether
it is actually "well-formed". This is just _my_ way.
Use this as you see fit. Replace the pieces that don't work for you. Let me know if you have "better" ways of doing things.
I'm not dogmatic about this; in fact, I'm always excited to see (other|easier|better|different) ways of doing things.
Sometimes those things will wrap around my brain better than what I'm doing now. This is progress! I may modify
this template in the future as I learn more and as new modules become available.

## Quick Start

You must have `go` and `curl` installed to fully run the demo. 

Build the binary:

```
go get
go build
```

Run the server:

```
./go-rest-api-template
```

The `demo.sh` script will walk through the steps below, if you don't want to cut and paste.

Create a policy for your user (I use the user/pass `demo`/`s3cr3t` below):

```
USER=demo
PASS=s3cr3t
echo "p, ${USER}, data, read" > ./config/policy.csv
```

Generate and store a password hash (the server will cache the hash):

```
curl -s -q -u ${USER}:${PASS} localhost:8000/generate
```

Get a token from the server:

```
token=$(curl -s -q -u ${USER}:${PASS} localhost:8000/auth | sed 's/^.*token":"//' | sed 's/"}.*$//')
```

Use the token to request "data":

```
curl -H "Content-type: application/json" -H "Token:$token" localhost:8000/data/${USER}
```

## Testing

In order to test properly, the `config/policy.csv` file must contain:

```
p, demo, data, read
```

Run all of the tests:

```
go test
```

Run a single test:

```
go test -run TestResourceHandler
```

## See Also

* [pete0emerson/go-cli-template](https://github.com/pete0emerson/go-cli-template) is a similar project for a well-formed CLI in Go.

## Future improvements

* ??? `/metrics` endpoint exposing metrics for ingestion into something like Prometheus
* `/ready` endpoint for a load balancer to check
* `/serve` authenticated POST endpoint to set desired serving state (to remove from load balancer even if the server is healthy)
* `/version` endpoint
* ??? Build scripts (`Makefile`? `/scripts/foo.sh`?)
* Use `docker-compose` to build a working local environment
  * ??? App
  * ??? Prometheus
  * Grafana
  * ??? Redis
  * Caching layer
* Minikube deployment
* Homebrew packaging
* Tie Casbin to Redis with [an adapter](https://github.com/mlsen/casbin-redis-adapter) for persistence instead of files