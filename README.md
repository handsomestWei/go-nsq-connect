# go-nsq-connect
golang消息中间件**nsq**的连接端，包含生产者和消费者。以及nsq消息生产和消费的demo

## nsq消息流
图源官方文档
<div align=center><img width="420" height="281" src="https://github.com/handsomestWei/go-nsq-connect/blob/main/nsq-design.gif" /></div>

## nsq组件介绍
+ **nsqlookupd**：中心管理服务，守护进程负责管理拓扑信息。客户端通过查询 nsqlookupd 来发现指定话题（topic）的生产者，并且 nsqd 节点广播话题（topic）和通道（channel）信息。
+ **nsqd**：一个守护进程，负责接收，排队，投递消息给客户端。对订阅了同一个topic，同一个channel的消费者使用负载均衡策略。
+ **nsqadmin**：一套WEB UI，用来汇集集群的实时统计，并执行不同的管理任务。

## nsq docker部署
nsq主要有三个组件: nsqlookupd, nsqd, nsqadmin。这三个组件都包含在nsqio/nsq镜像中
```
docker pull nsqio/nsq
```
### 部署nsqlookupd
```
docker run --name nsqlookupd -p 4160:4160 -p 4161:4161 -d nsqio/nsq /nsqlookupd
```

### 部署nsqd
```
docker run --name nsqd -p 4150:4150 -p 4151:4151 -d nsqio/nsq /nsqd --broadcast-address=172.16.21.11 --lookupd-tcp-address=172.16.21.11:4160 --data-path=/data
```

### 部署nsqadmin
```
docker run -d --name nsqadmin -p 4171:4171 nsqio/nsq /nsqadmin --lookupd-http-address=172.16.21.11:4161
```

### 控制台访问
```
http://172.16.21.11:4171
```

## 参考
[nsq连接](https://github.com/nsqio/go-nsq)   
[nsq连接说明](https://godoc.org/github.com/nsqio/go-nsq)   
[nsq官方文档](https://nsq.io/overview/quick_start.html)