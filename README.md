# cNET

[![GitHub license](https://img.shields.io/github/license/tomhjx/netcat.svg?style=popout-square)](https://github.com/tomhjx/netcat/blob/main/LICENSE) ![GitHub Workflow Status (Publish)](https://img.shields.io/github/actions/workflow/status/tomhjx/cnet/publish.yml) ![Docker Pulls](https://img.shields.io/docker/pulls/tomhjx/cnet)

便捷的网络探针，收集网络请求及响应过程信息

## 开始

* [Docker Image](https://hub.docker.com/r/tomhjx/cnet)

```bash
# 开发版本
docker run --rm tomhjx/cnet:develop -h

# 最新可用版本
docker run --rm tomhjx/cnet:main -h

# tag x.x.x 版本
docker run --rm tomhjx/cnet:x.x.x -h

```

* 编译源码

依赖 `go 1.21+`

```bash

go build -o cnet ./cmd/cnet/main.go

cnet -h

```

## 能力

* 探测协议
    * [x] IP
    * [x] Domain
    * [ ] TCP
    * [ ] UDP
    * [ ] Web Socket
    * [ ] Socket
    * [x] AMQP
    * [x] HTTP
    * [x] HTTPs

* 输出格式
    * [x] JSON

* 输出方式
    * [x] STDOUT
    * [x] SYSLOG
    * [ ] UDP
    * [ ] TCP
    * [ ] Unix Socket

* 接入监控
    * [x] Prometheus

## 使用
* [配置示例](./examples)
* [字段含义](./doc/fields.md)
* [查看日志](./doc/logging.md)

