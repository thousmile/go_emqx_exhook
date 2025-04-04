## 快速启动

[感谢@JetBrains 提供 Goland 支持![](./images/goland.svg)](https://www.jetbrains.com/?from=spring-admin-vue)

vim /etc/go_emqx_exhook/config.yaml

```yaml
appName: go_emqx_exhook
port: 16565

# mq类型: Rocketmq、Kafka、Rabbitmq、RabbitmqStream、Redis
mqType: Rocketmq

# emqx 主题
bridgeRule:
  topics:
    - "/#"

# rocketmq 配置，需要提前创建 主题
rocketmqConfig:
  nameServer:
    - 127.0.0.1:9876
  topic: emqx_exhook
  tag: exhook
  groupName: exhook
  #accessKey: exhook
  #secretKey: exhook


# rabbitmq 配置，需要提前创建 交换机 并且绑定队列
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


# rabbitmq stream 配置
rabbitmqStreamConfig:
  addresses:
    - rabbitmq-stream://guest:guest@127.0.0.1:5552
  # 主题不存在，就自动创建
  streamName: emqx_exhook
  # 发送者数量
  maxProducersPerClient: 2
  # x-max-age 支持 [ s, m, h ] 默认: 168h
  maxAge: 168h
  # x-max-length-bytes 支持 [ kb, mb, gb, tb ] 默认: "10gb"
  maxLengthBytes: 10gb
  # x-stream-max-segment-size-bytes 支持 [ kb, mb, gb, tb ] 默认: "1gb"
  maxSegmentSizeBytes: 1gb
  # 支持: "none", "gzip", "snappy", "lz4", "zstd", 默认: "none"
  compressionCodec: zstd
#  tls:
#    enable: true
#    tlsSkipVerify: true
#    caFile: /apps/server.cer.pem
#    certFile: /apps/client.cer.pem
#    keyFile: /apps/client.key.pem


# kafka 配置
kafkaConfig:
  addresses:
    - 127.0.0.1:9092
  # 主题不存在，就自动创建
  topic: emqx_exhook
  # 消息压缩类型 支持: "none", "gzip", "snappy", "lz4", "zstd", 默认: "none"
  compressionCodec: none
  # 分区数量，默认 -1 
  numPartitions: 3
  # 副本数量，默认 -1 
  replicationFactor: 3
  # 自定主题配置
  configEntries:
    retention.ms: 604800000

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


# redis 配置
redisConfig:
  addresses:
    - 127.0.0.1:6379
  # 主题不存在，就自动创建
  streamName: emqx_exhook
  # 主题最大的消息数量，超出会自动移除最开始的消息，-1表示没有限制
  streamMaxLen: 100000
  db: 0
#  username: redis123
#  password: redis123456
#  masterName: mymaster
#  sentinelUsername: sentinel123456
#  sentinelPassword: sentinel123456



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

grpc server 支持 tls
/etc/go_emqx_exhook/config.yaml 添加以下配置

```shell
appName: go_emqx_exhook
port: 16565
# grpc 支持 tls
tls:
  enable: true
  caFile: certs/ca/ca.crt
  certFile: certs/server/server.crt
  keyFile: certs/server/server.key
```

<span style="color:red;"> Emqx > ExHook > URL: 必须以 https 开始，如: https://127.0.0.1:16565 </span>

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
docker build -t go_emqx_exhook:1.8 ./


## 运行docker容器
docker run -d --name go_emqx_exhook -p 16565:16565 --restart=always go_emqx_exhook:1.8


## 指定配置文件
docker run -d --name go_emqx_exhook -p 16565:16565 \
  -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml \ 
  -v /etc/localtime:/etc/localtime:ro \ 
  --restart=always thousmile/go_emqx_exhook:1.8

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

[kafka-generate-ssl-automatic.sh](https://github.com/confluentinc/confluent-platform-security-tools)

### jks 转换 pem

```shell

keytool -importkeystore -srckeystore kafka.truststore.jks -destkeystore server.p12 -deststoretype PKCS12

openssl pkcs12 -in server.p12 -nokeys -out server.cer.pem

keytool -importkeystore -srckeystore kafka.keystore.jks -destkeystore client.p12 -deststoretype PKCS12

openssl pkcs12 -in client.p12 -nokeys -out client.cer.pem

openssl pkcs12 -in client.p12 -nodes -nocerts -out client.key.pem

```

