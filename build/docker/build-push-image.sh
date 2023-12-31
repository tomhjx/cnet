#!/bin/bash
root=$(dirname $(dirname $(dirname "$0")))
pwd=$(dirname "$0")
tag=$1

echo $root
echo $pwd
echo $tag
# docker run --rm -w /app -v ${root}:/app -v ${pwd}:/out golang:1.21-alpine3.18 go build -o /out/cnet ./cmd/cnet/main.go
# ls -anl ${pwd}/cnet
docker buildx build --platform linux/amd64,linux/arm64 ${root} -f ${pwd}/Dockerfile -t ${tag} --push