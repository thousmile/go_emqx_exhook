package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go_emqx_exhook/channelx"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"go_emqx_exhook/impl"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

func main() {
	appConf := conf.Config

	rmqConf := conf.Config.RocketmqConfig
	rmqRule := conf.Config.BridgeRule
	rmqQueue := conf.Config.Queue
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

	// 批量消息队列中
	aggr := channelx.NewAggregator[*exhook.Message](
		func(messages []*exhook.Message) error {
			targetMessages := make([]*primitive.Message, len(messages))
			for idx, sourceMessage := range messages {
				targetMessage := &primitive.Message{
					Topic: rmqRule.TargetTopic,
					Body:  sourceMessage.Payload,
				}
				targetMessage.WithKeys([]string{sourceMessage.Id})
				targetMessage.WithTag(rmqRule.TargetTag)
				for key, val := range sourceMessage.GetHeaders() {
					targetMessage.WithProperty(key, val)
				}
				targetMessage.WithProperty("sourceId", sourceMessage.Id)
				targetMessage.WithProperty("sourceTopic", sourceMessage.Topic)
				targetMessage.WithProperty("sourceNode", sourceMessage.Node)
				targetMessage.WithProperty("sourceFrom", sourceMessage.From)
				targetMessage.WithProperty("sourceQos", strconv.Itoa(int(sourceMessage.Qos)))
				targetMessage.WithProperty("sourceTimestamp",
					strconv.FormatInt(int64(sourceMessage.Timestamp), 10))
				targetMessages[idx] = targetMessage
			}
			_, err = rmqProducer.SendSync(context.Background(), targetMessages...)
			if err != nil {
				log.Printf("batch send message [%d] error: %s\n", len(targetMessages), err)
			}
			return nil
		},
		func(option channelx.AggregatorOption[*exhook.Message]) channelx.AggregatorOption[*exhook.Message] {
			option.BatchSize = rmqQueue.BatchSize
			option.Workers = rmqQueue.Workers
			option.ChannelBufferSize = option.BatchSize * 2
			option.LingerTime = time.Duration(rmqQueue.LingerTime) * time.Second
			option.Logger = log.Default()
			log.Printf("channelx option : %v \n", option)
			return option
		},
	)
	aggr.Start()
	defer aggr.SafeStop()

	srv := grpc.NewServer()

	// 注册 emqx 的 exhook grpc 服务
	exhook.RegisterHookProviderServer(srv, &impl.HookProviderServerImpl{
		SourceTopics: rmqRule.SourceTopics,
		// 接收到 emqx 的消息后，立即发送到 队列中。
		Receive: func(msg *exhook.Message) {
			aggr.TryEnqueue(msg)
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
	log.Printf("%s => grpc server listen port : %d \n", appConf.AppName, appConf.Port)
	_ = srv.Serve(lis)
}
