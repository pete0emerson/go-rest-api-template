#!/bin/bash

USER=demo
PASS=s3cr3t

echo "## Generating Casbin policy file in ./config/policy.csv"
echo "p, ${USER}, data, read" > ./config/policy.csv
echo

echo "## Generating a password (the server will cache it)"
curl -s -q -u ${USER}:${PASS} localhost:8000/generate
echo

echo "## Getting a token from the server"
token=$(curl -s -q -u ${USER}:${PASS} localhost:8000/auth | cut -d , -f 3 | cut -d : -f 2 | sed 's/["}]//g')
echo "## Got a token from the server"
if [[ "$token" == "null" ]] ; then
	echo "## Not authorized. Quitting."
	exit 1
fi
echo

echo "## Using the token to request /data/${USER}"
curl -H "Content-type: application/json" -H "Token:$token" localhost:8000/data/${USER}
echo
