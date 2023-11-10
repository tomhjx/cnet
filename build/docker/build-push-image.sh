#!/bin/bash
root=$(dirname $(dirname $(dirname "$0")))
pwd=$(dirname "$0")
tag=$1

echo $pwd
echo $dir
echo $tag
docker run --rm -v ${root}:/app -v ${pwd}:/out golang:1.21-alpine3.18 go build -o /out/cnet /app/cmd/cnet/main.go
ls -anl ${pwd}/cnet
docker buildx build --platform linux/amd64,linux/arm64 ${pwd} -t $tag --push