## 快速运行

vim /etc/go_emqx_exhook/config.yaml

```yaml
appName: go_emqx_exhook
port: 16565

# Rocketmq、Rabbitmq、Kafka
mqType: Rocketmq

# emqx 主题
bridgeRule:
  topics:
    - "/#"

# rocketmq 配置，需要提前创建 主题
rocketmqConfig:
  nameServer:
    - 192.168.0.188:9876
  topic: emqx_exhook
  tag: exhook
  groupName: exhook
  #accessKey: exhook
  #secretKey: exhook


# rabbitmq 配置，需要提前创建 交换机 并且绑定队列
rabbitmqConfig:
  addresses:
    - amqp://rabbit:mht123456@192.168.0.188:5672
  exchangeName: emqx_exhook
  routingKeys: exhook


# kafka 配置，需要提前创建 主题
kafkaConfig:
  addresses:
    - 192.168.0.188:9092
  topic: emqx_exhook


# 发送方式 queue or direct ，默认 queue
# 注: rabbitmq 不支持队列发送
sendMethod: queue

queue:
  batchSize: 100
  workers: 2
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

## 本地运行

[根据自己的操作系统，下载相应的 可执行文件 ](https://github.com/thousmile/go_emqx_exhook/releases)
解压缩后，在 可执行文件 同级目录下，新建 config.yaml 配置文件

## 在 EMQX Dashboard > ExHook

![](./images/20230728154744.png)

### 二次开发

```shell
proto 生成 go 代码
protoc --go_out=. --go-grpc_out=. proto/*.proto


golang 打包可执行文件
goreleaser --snapshot --skip-publish --clean


构建docker镜像
docker build -t go_emqx_exhook:1.3 ./


运行docker容器
docker run -d --name go_emqx_exhook -p 16565:16565 --restart=always go_emqx_exhook:1.3


## 指定配置文件
docker run -d --name go_emqx_exhook -p 16565:16565 \
  -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml \ 
  -v /etc/localtime:/etc/localtime:ro \ 
  --restart=always thousmile/go_emqx_exhook:1.3

```

### 消费者 获取mqtt的属性或者Header

| 属性              | 描述                     |
|-----------------|------------------------|
| sourceId        | mqtt 消息ID              |
| sourceTopic     | mqtt 主题                |
| sourceNode      | emqx 节点名称              |
| sourceFrom      | emqx 来自哪个mqtt客户端       |
| sourceQos       | mqtt qos               |
| sourceTimestamp | 消息时间戳                  |
| protocol        | 此消息协议(emqx默认Header)    |
| peerhost        | 此消息生产者IP(emqx默认Header) |


Rabbitmq:
![](./images/20231207160607.png)


Kafka:
![](./images/20231207164403.png)

