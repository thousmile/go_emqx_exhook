package main

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"go_emqx_exhook/impl"
	"go_emqx_exhook/provider"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

func main() {
	appConf := conf.Config
	rule := appConf.BridgeRule
	// 创建一个消息提供者
	var msgProvider provider.MessageProvider

	if strings.EqualFold(appConf.MqType, "Rabbitmq") || strings.EqualFold(appConf.MqType, "rabbitmq") {
		rabbitmqProvider := provider.BuildRabbitmqMessageProvider(appConf.RabbitmqConfig, rule.TargetTopic)
		msgProvider = rabbitmqProvider
	} else {
		rocketmqProvider := provider.BuildRocketmqMessageProvider(appConf.RocketmqConfig, rule.TargetTopic)
		defer func(p rocketmq.Producer) {
			err := p.Shutdown()
			if err != nil {
				log.Fatal(err)
			}
		}(rocketmqProvider.RmqProducer)
		msgProvider = rocketmqProvider
	}

	ch := make(chan *exhook.Message, appConf.ChanBufferSize)

	// 发送方式“ queue or direct ”
	if appConf.SendMethod == "queue" {
		go Queue(msgProvider, ch)
	} else {
		go Direct(msgProvider, ch)
	}

	srv := grpc.NewServer()
	// 注册 emqx 的 exhook grpc 服务
	exhook.RegisterHookProviderServer(srv, &impl.HookProviderServerImpl{
		SourceTopics: rule.SourceTopics,
		Receive:      ch,
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
