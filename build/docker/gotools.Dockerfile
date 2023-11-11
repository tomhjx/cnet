FROM golang:1.21-alpine3.18
RUN apk add --update graphviz && \
    rm -rf /var/cache/apk/*