#!/bin/bash

USER=demo
PASS=s3cr3t

echo "## Generating a password (the server will cache it)"
curl -s -q -u ${USER}:${PASS} localhost:8000/generate
echo

echo "## Getting a token from the server"
token=$(curl -s -q -u ${USER}:${PASS} localhost:8000/auth | sed 's/^.*token":"//' | sed 's/"}.*$//')
echo

echo "## Got a token from the server"
if [[ "$token" == "null" ]] ; then
	echo "## Not authorized. Quitting."
	exit 1
fi
echo

echo "## Using the token to request /data/${USER}"
curl -H "Content-type: application/json" -H "Token:$token" localhost:8000/data/${USER}
echo

echo "## Testing redis"
curl -s -q localhost:8000/redis
curl -s -q localhost:8000/redis
curl -s -q localhost:8000/redis
