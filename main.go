package main

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"go_emqx_exhook/impl"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	appConf := conf.Config

	rmqConf := appConf.RocketmqConfig
	rmqRule := appConf.BridgeRule
	rmqProducer, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(rmqConf.NameServer)),
		producer.WithRetry(2),
		producer.WithGroupName("exhook"),
		producer.WithSendMsgTimeout(time.Second*1),
	)
	err := rmqProducer.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer func(rmqProducer rocketmq.Producer) {
		err := rmqProducer.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}(rmqProducer)

	srv := grpc.NewServer()

	ch := make(chan *exhook.Message, 10000)
	// 发送方式“ queue or direct ”
	if appConf.SendMethod == "queue" {
		go Queue(rmqProducer, ch)
	} else {
		go Direct(rmqProducer, ch)
	}

	// 注册 emqx 的 exhook grpc 服务
	exhook.RegisterHookProviderServer(srv, &impl.HookProviderServerImpl{
		SourceTopics: rmqRule.SourceTopics,
		// 接收到 emqx 的消息后，立即发送到 队列中。
		Receive: func(msg *exhook.Message) {
			ch <- msg
		},
	})

	// 监听 指定端口的 tcp 连接
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", appConf.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(lis net.Listener) {
		err := lis.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(lis)
	log.Printf("%s [%s] => grpc server listen port : %d \n", appConf.AppName, appConf.SendMethod, appConf.Port)
	_ = srv.Serve(lis)
}
