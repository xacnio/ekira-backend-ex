#!/bin/bash
DIRNAME=$(dirname "$0")
DIR=$(realpath "$DIRNAME")
USERNAME=xacnio
IMAGE=ekira-backend
BASE=$(realpath "$DIR/..")
VERSION=$(cat $BASE/VERSION | tr --delete '\r' | tr --delete '\n')
cd "$BASE"
docker build -t $USERNAME/$IMAGE:latest .
docker tag $USERNAME/$IMAGE:latest $USERNAME/$IMAGE:$VERSION
