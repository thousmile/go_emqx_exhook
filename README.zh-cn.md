## 快速启动

vim /etc/go_emqx_exhook/config.yaml

```yaml
appName: go_emqx_exhook
port: 16565

# mq类型: Rocketmq、Kafka、Rabbitmq、Redis
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
#  sasl:
#    enable: true
#    user: admin
#    password: 123456


# redis 配置，无需创建 stream 
redisConfig:
  addresses:
    - 127.0.0.1:6379
  streamName: emqx_exhook
  db: 0
  username: redis123
  password: redis123456
  masterName: mymaster
  sentinelUsername: sentinel123456
  sentinelPassword: sentinel123456


# 发送方式 queue 或者 direct ，默认 queue
# queue: 收到消息后，转入队列，当队列内的消息数量等于阈值，批量发送到mq中
# direct: 收到消息后，立即发送到mq中
# 注: rabbitmq 和 redis 不支持队列发送
sendMethod: queue

# 队列的配置， batchSize 和 lingerTime 只要满足一个，就将消息批量发送到mq中
queue:
  # 当消息数量达到100条是，批量发送到mq中
  batchSize: 100
  workers: 2
  # 收到消息后，无论队列中的消息数量是否满足，都会在1秒内发送出去。
  lingerTime: 1

```

```shell
docker run -d --name go_emqx_exhook -p 16565:16565 \
  -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml \
  -v /etc/localtime:/etc/localtime:ro \
  --restart=always thousmile/go_emqx_exhook:1.5
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
    image: thousmile/go_emqx_exhook:1.5
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
## proto 生成 go 代码
protoc --go_out=. --go-grpc_out=. proto/*.proto


## golang 打包可执行文件
goreleaser --snapshot --skip-publish --clean


## 构建docker镜像
docker build -t go_emqx_exhook:1.5 ./


## 运行docker容器
docker run -d --name go_emqx_exhook -p 16565:16565 --restart=always go_emqx_exhook:1.5


## 指定配置文件
docker run -d --name go_emqx_exhook -p 16565:16565 \
  -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml \ 
  -v /etc/localtime:/etc/localtime:ro \ 
  --restart=always thousmile/go_emqx_exhook:1.5

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



Redis:
![](./images/20240201103222.png)

Rabbitmq:
![](./images/20231207160607.png)

Kafka:
![](./images/20231207164403.png)

