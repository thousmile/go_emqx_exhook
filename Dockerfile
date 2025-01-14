FROM golang:bookworm AS builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -o /build/app

FROM debian:stable-slim
WORKDIR /apps/
COPY --from=builder --chown=1001:root /build/app /apps/app
COPY --from=builder --chown=1001:root /build/config.yaml /apps/config.yaml
EXPOSE 16565
CMD ["./app"]
