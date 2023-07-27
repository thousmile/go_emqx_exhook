### proto 生成 go 代码
```shell

protoc --go_out=. --go-grpc_out=. proto/*.proto

protoc -I=proto --go_out=. --go-grpc_out=. proto/*.proto

```

### golang 编译
```shell
## 打包可执行文件
goreleaser build --single-target


```

### 下载 Linux 可执行文件。

vim Dockerfile
```shell
FROM alpine:3.18.2
MAINTAINER Wang Chen Chen<932560435@qq.com>
ENV VERSION 1.0
# 在容器根目录 创建一个 apps 目录
WORKDIR /apps
# 拷贝当前目录下 go_docker_demo1 可以执行文件
COPY go_emqx_exhook /apps/golang_app
# 拷贝配置文件到容器中
COPY config.yaml /apps/config.yaml
# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
# 设置编码
ENV LANG C.UTF-8
# 暴露端口
EXPOSE 16565
# 运行golang程序的命令
ENTRYPOINT ["/apps/golang_app"]
```

vim config.yaml
```yaml
appName: go_emqx_exhook
port: 16565

bridgeRule:
  sourceTopics:
    - "/#"
  targetTopic: emqx_msg_bridge
  targetTag: emqx

rocketmqConfig:
  nameServer:
    - 127.0.0.1:9876
```

### docker build
```shell
docker build -t go_emqx_exhook:1.1 ./
```

## docker 运行
```shell
docker run -d --name go_emqx_exhook -p 16565:16565 --restart=always go_emqx_exhook:1.1

## 指定配置文件
docker run -d --name go_emqx_exhook -p 16565:16565 -v /etc/config.yaml:/apps/config.yaml --restart=always go_emqx_exhook:1.1
```