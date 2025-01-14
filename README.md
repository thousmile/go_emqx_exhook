## Quick Start

### [中文文档](README.zh-cn.md)

[Thanks to @JetBrains for providing Goland support ![](./images/goland.svg)](https://www.jetbrains.com/?from=spring-admin-vue)

vim /etc/go_emqx_exhook/config.yaml

```yaml
appName: go_emqx_exhook
port: 16565

# Rocketmq or Rabbitmq or RabbitmqStream or Kafka or Redis
mqType: Rocketmq

# grpc server tls
tls:
  enable: true
  caFile: certs/ca/ca.crt
  certFile: certs/server/server.crt
  keyFile: certs/server/server.key

# emqx topic
bridgeRule:
  topics:
    - "/#"

# rocketmq configuration, you need to create a topic in advance
rocketmqConfig:
  nameServer:
    - 127.0.0.1:9876
  topic: emqx_exhook
  tag: exhook
  groupName: exhook
  #accessKey: exhook
  #secretKey: exhook


# rabbitmq configuration, you need to create a switch in advance and bind a queue
rabbitmqConfig:
  addresses:
    - amqp://guest:guest@127.0.0.1:5672
  exchangeName: emqx_exhook
  routingKeys: exhook
#  tls:
#    enable: true
#    tlsSkipVerify: true
#    caFile: /apps/server.cer.pem
#    certFile: /apps/client.cer.pem
#    keyFile: /apps/client.key.pem


# rabbitmq stream configuration, you need to create a stream in advance
rabbitmqStreamConfig:
  addresses:
    - rabbitmq-stream://guest:guest@127.0.0.1:5552
  streamName: emqx_exhook
  # number of senders
  maxProducersPerClient: 2
#  tls:
#    enable: true
#    tlsSkipVerify: true
#    caFile: /apps/server.cer.pem
#    certFile: /apps/client.cer.pem
#    keyFile: /apps/client.key.pem


# kafka configuration, you need to create a topic in advance
kafkaConfig:
  addresses:
    - 127.0.0.1:9092
  topic: emqx_exhook
  # compression type support: "none", "gzip", "snappy", "lz4", "zstd", default: "none"
  compressionCodec: none
#  sasl:
#    enable: true
#    user: admin
#    password: admin123456
#  tls:
#    enable: true
#    tlsSkipVerify: true
#    caFile: /apps/server.cer.pem
#    certFile: /apps/client.cer.pem
#    keyFile: /apps/client.key.pem


# redis configuration
redisConfig:
  addresses:
    - 127.0.0.1:6379
  streamName: emqx_exhook
  # The maximum number of messages in a topic. If the number exceeds the limit, the first message will be automatically removed. -1 means no limit
  streamMaxLen: 100000
  db: 0
#  username: redis123
#  password: redis123456
#  masterName: mymaster
#  sentinelUsername: sentinel123456
#  sentinelPassword: sentinel123456


# message send method "queue or direct", default: queue
# queue: after receiving the message, enter the queue and send it in batch when the queue conditions are met.
# direct: send immediately after receiving the message
# info: Rabbitmq and Redis do not support queue
sendMethod: queue


# queue configuration batchSize and lingerTime only satisfy one of them
queue:
  # when the number of messages in the queue reaches 100, batch send
  batchSize: 100
  workers: 2
  # after receiving the message, regardless of whether the number of messages in the queue is satisfied, it will be sent within 1 second.
  lingerTime: 1

```

grpc server support tls
append to /etc/go_emqx_exhook/config.yaml 

```shell
appName: go_emqx_exhook
port: 16565
# grpc server tls
tls:
  enable: true
  caFile: certs/ca/ca.crt
  certFile: certs/server/server.crt
  keyFile: certs/server/server.key
```

<span style="color:red;"> Emqx > ExHook > URL: Must start with https for example : https://127.0.0.1:16565 </span>

```shell
docker run -d --name go_emqx_exhook -p 16565:16565 \
  -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml \
  -v /etc/localtime:/etc/localtime:ro \
  --restart=always thousmile/go_emqx_exhook:1.8
```

vim docker-compose.yml

```yaml
version: '3'

networks:
  app-net1:
    ipam:
      config:
        - subnet: 172.19.0.0/16
          gateway: 172.19.0.1

services:
  go_emqx_exhook:
    image: thousmile/go_emqx_exhook:1.8
    container_name: go_emqx_exhook
    ports:
      - "16565:16565"
    volumes:
      - /etc/go_emqx_exhook/config.yaml:/apps/config.yaml
      - /etc/localtime:/etc/localtime:ro
    privileged: true
    restart: always
    networks:
      app-net1:
    deploy:
      resources:
        limits:
          memory: 258m

```

```shell
docker compose up -d go_emqx_exhook
```

## binary

[download binary file](https://github.com/thousmile/go_emqx_exhook/releases)
after decompression, create a new config.yaml configuration file in the same directory as the binary file

## EMQX Dashboard > ExHook

![](./images/20231207180925.png)

### based on this project

```shell
# proto generate golang code
protoc --go_out=. --go-grpc_out=. proto/*.proto


# package binary
goreleaser --snapshot --skip-publish --clean
goreleaser release --skip=publish --clean

# build docker image
docker build -t go_emqx_exhook:1.8 ./


# run docker container
docker run -d --name go_emqx_exhook -p 16565:16565 --restart=always go_emqx_exhook:1.8


## custom configuration file
docker run -d --name go_emqx_exhook -p 16565:16565 \
  -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml \ 
  -v /etc/localtime:/etc/localtime:ro \ 
  --restart=always thousmile/go_emqx_exhook:1.8

```

### consumer get mqtt attributes or Header

| attributes name | description              |
|-----------------|--------------------------|
| sourceId        | mqtt message Id          |
| sourceTopic     | mqtt topic               |
| sourceNode      | emqx node name           |
| sourceFrom      | emqx from mqtt client id |
| sourceQos       | mqtt qos                 |
| sourceTimestamp | message timestamp        |
| protocol        | message protocol         |
| peerhost        | producer ip              |



Redis:
![](./images/20240201103222.png)

Rabbitmq:
![](./images/20231207160607.png)

Kafka:
![](./images/20231207164403.png)



[kafka-generate-ssl-automatic.sh](https://github.com/confluentinc/confluent-platform-security-tools)

### jks convert pem

```shell

keytool -importkeystore -srckeystore kafka.truststore.jks -destkeystore server.p12 -deststoretype PKCS12

openssl pkcs12 -in server.p12 -nokeys -out server.cer.pem

keytool -importkeystore -srckeystore kafka.keystore.jks -destkeystore client.p12 -deststoretype PKCS12

openssl pkcs12 -in client.p12 -nokeys -out client.cer.pem

openssl pkcs12 -in client.p12 -nodes -nocerts -out client.key.pem

```

