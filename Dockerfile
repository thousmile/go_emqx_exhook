FROM --platform=$TARGETPLATFORM alpine:3.18.6
ARG TARGETPLATFORM
ARG TARGETARCH
MAINTAINER Wang Chen Chen<932560435@qq.com>
ENV VERSION 1.0
WORKDIR /apps
COPY dist/go_emqx_exhook_linux_${TARGETARCH}/go_emqx_exhook /apps/golang_app
COPY config.yaml /apps/config.yaml
RUN chown 1001 /apps
RUN chmod "g+rwX" /apps
RUN chown 1001:root /apps
ENV LANG C.UTF-8
EXPOSE 16565
ENTRYPOINT /apps/golang_app