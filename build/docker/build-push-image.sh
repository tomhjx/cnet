#!/bin/bash
root=`dirname $0`
tag=$1

echo $root
echo $tag
docker run --rm -v /Users/tom/Work/project/github.com/tomhjx/cnet:/app -v ${root}:/out golang:1.21-alpine3.18 go build -o /out/cnet /app/cmd/cnet/main.go
ls -anl ${root}/cnet
docker buildx build --platform linux/amd64,linux/arm64 ${root} -t $tag --push