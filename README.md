## 快速运行

vim /etc/go_emqx_exhook/config.yaml

```yaml
appName: go_emqx_exhook
port: 16565
chanBufferSize: 10240

bridgeRule:
  ## emqx 的主题，可以多个
  sourceTopics:
    - "/test/#"
    - "/hello/+"
  ## rocketmq 的 主题
  targetTopic: emqx_msg_bridge
  targetTag: emqx

## rocketmq name server
rocketmqConfig:
  nameServer:
    - 127.0.0.1:9876

## 发送方式 queue or direct ，默认 queue
## direct: 收到消息后，立即转发到 rocketmq 中
## queue: 收到消息后，消息进入队列，当满足任意条件(1.收到100条消息，2.收到消息的时间大于1秒)，批量转发到 rocketmq。 
sendMethod: queue

## 队列信息
queue:
  batchSize: 100
  workers: 2
  lingerTime: 1

```

docker run -d --name go_emqx_exhook -p 16565:16565 -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml --restart=always thousmile/go_emqx_exhook:1.1

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
docker build -t go_emqx_exhook:1.1 ./


运行docker容器
docker run -d --name go_emqx_exhook -p 16565:16565 --restart=always go_emqx_exhook:1.1

## 指定配置文件
docker run -d --name go_emqx_exhook -p 16565:16565 -v /etc/go_emqx_exhook/config.yaml:/apps/config.yaml --restart=always go_emqx_exhook:1.1

```
