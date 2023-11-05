# cNET

便捷的网络探针，收集网络请求及响应过程信息

[![GitHub license](https://img.shields.io/github/license/tomhjx/netcat.svg?style=popout-square)](https://github.com/tomhjx/netcat/blob/main/LICENSE)

## 能力

* 统计各阶段耗时
* 支持prometheus接入
* 支持多个网络协议
* 支持主流输出协议


### 支持协议

* [x] IP
* [x] Domain
* [ ] TCP
* [ ] UDP
* [ ] Web Socket
* [ ] Socket
* [x] AMQP
* [x] HTTP
* [x] HTTPs


### 输出格式

* [x] JSON

### 输出源
* [x] STDOUT
* [x] SYSLOG
* [ ] UDP
* [ ] TCP
* [ ] Unix Socket

## 使用

```bash

cnet -h

```


## 字段

### 内容字段

字段名           | 含义
----------------|-----
id              | 单次上报处理的ID，用于跟踪处理片段（单次请求，或者批处理）信息
jid             | 上报任务的ID，用于跟踪任务情况
tid             | 上报子任务的ID，用于跟踪子任务情况，`jid`一样的子任务即由同一个主任务调度生成
cid             | 用于标记上报客户端，可自定义
url             | url地址
method          | 请求方法，GET/POST/PUT/DELETE...
connected_via   | 连接过程中相关的标记
headers         | 响应头信息
body            | 响应正文
status          | 响应状态描述
status_code     | 响应状态码 



### 耗时统计字段，单位秒

* http/https

```
+----------------->--------------->--------------->--------------->--------------------->----------------->---------------+
|    DNS Lookup   | TCP Handshake | SSL Handshake | Server Accept |  Server Processing  | Header Transfer | Body Transfer |
∧-----------------∧---------------∧---------------∧---------------∧---------------------∧-----------------+---------------∧
|name_lookup_time |   tcp_time    |    ssl_time   |               | server_process_time |       content_transfer_time     |
+-----------------+---------------+---------------+               +---------------------+---------------------------------+
|       connect_time              |               |               |                     |                                 |
+---------------------------------+               |               |                     |                                 |
|                   app_connect_time              |               |                     |                                 |
+-------------------------------------------------+               |                     |                                 |
|                       pre_transfer_time         |               |                     |                                 |
+-------------------------------------------------+---------------+                     |                                 |
|                                                 |           ttfb_time                 |                                 |
|                                                 +-------------------------------------+                                 |
|                                           start_transfer_time                         |                                 |
+---------------------------------------------------------------------------------------+                                 |
|                                                 total_time                                                              |
+-------------------------------------------------------------------------------------------------------------------------+


```

* host (ping)

```

+--------+----------+------------+--------+
| Client | Router(1)| Router(...)| Server |
+-------->---------->------------>--------+
|   total_time                            |
+-----------------------------------------+


```