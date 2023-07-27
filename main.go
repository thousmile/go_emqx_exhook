package main

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
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

	rmqConf := conf.Config.RocketmqConfig
	rmqProducer, _ := rocketmq.NewProducer(
		producer.WithNameServer(rmqConf.NameServer),
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

	rule := conf.Config.BridgeRule

	srv := grpc.NewServer()
	exhook.RegisterHookProviderServer(srv, &impl.HookProviderServerImpl{
		Producer:     rmqProducer,
		TargetTopic:  rule.TargetTopic,
		TargetTag:    rule.TargetTag,
		SourceTopics: rule.SourceTopics,
	})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", appConf.Port))
	defer func(lis net.Listener) {
		err := lis.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(lis)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("%s => grpc server listen port : %d \n", appConf.AppName, appConf.Port)
	_ = srv.Serve(lis)
}
