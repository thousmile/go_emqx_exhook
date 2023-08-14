FROM alpine:3.18.2
MAINTAINER Wang Chen Chen<932560435@qq.com>
ENV VERSION 1.0
# 在容器根目录 创建一个 apps 目录
WORKDIR /apps
# 拷贝当前目录下 go_emqx_exhook 可以执行文件
COPY dist/go_emqx_exhook_linux_amd64_v1/go_emqx_exhook /apps/golang_app
# 拷贝配置文件到容器中
COPY config.yaml /apps/config.yaml
# 设置编码
ENV LANG C.UTF-8
# 暴露端口
EXPOSE 16565
# 运行golang程序的命令
ENTRYPOINT /apps/golang_app