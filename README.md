# cnet
cNET

## 能力

* 统计各阶段耗时
* 支持多个网络协议
* 支持主流输出协议


### 支持协议

* [ ] IP
* [ ] TCP
* [ ] UDP
* [ ] Web Socket
* [ ] Socket
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
jid             | 上报任务的ID，用于跟踪任务情况
tid             | 上报子任务的ID，用于跟踪子任务情况，`jid`一样的子任务即由同一个主任务调度生成
cid             | 用于标记上报客户端，可自定义
url             | url地址
method          | 请求方法，GET/POST/PUT/DELETE...
connected_via   | 连接过程中相关的标记
headers         | 响应头信息
body            | 响应正文



### 耗时统计字段，单位毫秒

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
+                                                 +-------------------------------------+                                 |
|                                           start_transfer_time                         |                                 |
+---------------------------------------------------------------------------------------+                                 |
|                                                 total_time                                                              |
+-------------------------------------------------------------------------------------------------------------------------+



```