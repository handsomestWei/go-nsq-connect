## nsq主要有三个组件: nsqlookupd, nsqd, nsqadmin。这三个组件都包含在nsqio/nsq镜像中
docker pull nsqio/nsq

## 部署nsqlookupd
docker run --name nsqlookupd -p 4160:4160 -p 4161:4161 -d nsqio/nsq /nsqlookupd

## 部署nsqd
docker run --name nsqd -p 4150:4150 -p 4151:4151 -d nsqio/nsq /nsqd --broadcast-address=172.16.21.11 --lookupd-tcp-address=172.16.21.11:4160 --data-path=/data

## 部署nsqadmin
docker run -d --name nsqadmin -p 4171:4171 nsqio/nsq /nsqadmin --lookupd-http-address=172.16.21.11:4161