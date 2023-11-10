#!/bin/bash
root=`dirname $0`
tag=$1

echo $root
echo $tag
docker run -it --rm -w /app -v /Users/tom/Work/project/github.com/tomhjx/cnet:/app -v ${root}:/out golang:1.21-alpine3.18 go build -o /out/cnet ./cmd/cnet/main.go
ls -anl /out/cnet
docker buildx build --platform linux/amd64,linux/arm64 ${root} -t $tag --push