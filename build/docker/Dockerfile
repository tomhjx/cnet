FROM golang:1.21-alpine3.18 as builder
WORKDIR /app
COPY . /app
RUN mkdir /out && \
    go build -o /out/cnet ./cmd/cnet/main.go && \
    go version -m /out/cnet && \
    chmod +x /out/cnet

FROM alpine:3.18.4
COPY --from=builder ./cnet /usr/local/bin/cnet
ENTRYPOINT [ "cnet" ]