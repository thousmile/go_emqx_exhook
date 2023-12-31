## Quick Start

### [中文文档](README.zh-cn.md)

vim /etc/go_emqx_exhook/config.yaml

```yaml
appName: go_emqx_exhook
port: 16565

# Rocketmq or Rabbitmq or Kafka
mqType: Rocketmq

# emqx topic
bridgeRule:
  topics:
    - "/#"

# rocketmq configuration, you need to create a topic in advance
rocketmqConfig:
  nameServer:
    - 192.168.0.188:9876
  topic: emqx_exhook
  tag: exhook
  groupName: exhook
  #accessKey: exhook
  #secretKey: exhook


# rabbitmq configuration, you need to create a switch in advance and bind a queue
rabbitmqConfig:
  addresses:
    - amqp://rabbit:mht123456@192.168.0.188:5672
  exchangeName: emqx_exhook
  routingKeys: exhook


# kafka configuration, you need to create a topic in advance
kafkaConfig:
  addresses:
    - 192.168.0.188:9092
  topic: emqx_exhook


# message send method "queue or direct", default: queue
# queue: after receiving the message, enter the queue and send it in batch when the queue conditions are met.
# direct: send immediately after receiving the message
# info: rabbitmq queue send is not supported
sendMethod: queue


# queue configuration batchSize and lingerTime only satisfy one of them
queue:
  # when the number of messages in the queue reaches 100, batch send
  batchSize: 100
  workers: 2
  # after receiving the message, regardless of whether the number of messages in the queue is satisfied, it will be sent within 1 second.
  lingerTime: 1

```

```shell
docker run -d --name go_emqx_exhook -p 16565:16565 \
  -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml \
  -v /etc/localtime:/etc/localtime:ro \
  --restart=always thousmile/go_emqx_exhook:1.3
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
    image: thousmile/go_emqx_exhook:1.3
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


# build docker image
docker build -t go_emqx_exhook:1.3 ./


# run docker container
docker run -d --name go_emqx_exhook -p 16565:16565 --restart=always go_emqx_exhook:1.3


## custom configuration file
docker run -d --name go_emqx_exhook -p 16565:16565 \
  -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml \ 
  -v /etc/localtime:/etc/localtime:ro \ 
  --restart=always thousmile/go_emqx_exhook:1.3

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



Rabbitmq:
![](./images/20231207160607.png)

Kafka:
![](./images/20231207164403.png)

