```shell

docker run --rm --privileged tonistiigi/binfmt:latest --install all

docker buildx create --name mybuilder --driver docker-container

docker buildx use mybuilder

docker buildx build --platform linux/amd64,linux/arm64,linux/386 -t thousmile/go_emqx_exhook:1.7 . --push

```