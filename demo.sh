#!/bin/bash

USER=demo
PASS=s3cr3t

# Create the Casbin policy file
echo "p, ${USER}, data, read" > ./config/policy.csv

# Generate a password (the server will cache it)
hash=$(curl -s -q -u ${USER}:${PASS} localhost:8000/generate | jq -r '.hash')
echo "Got a password hash of $hash"
if [[ "$hash" == "null" ]] ; then
	echo "No hash generated."
	exit 1
fi

# Get a token from the server
token=$(curl -s -q -u ${USER}:${PASS} localhost:8000/auth | jq -r '.token')
echo "Got a token of $token"
if [[ "$token" == "null" ]] ; then
	echo "Not authorized. Quitting."
	exit 1
fi
echo "Requesting /data/${USER} ..."

# Use the token to request "data"
curl -H "Content-type: application/json" -H "Token:$token" localhost:8000/data/${USER}
