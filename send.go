package main

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go_emqx_exhook/channelx"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"log"
	"strconv"
	"time"
)

// Queue 使用队列
func Queue(rmqProducer rocketmq.Producer, ch chan *exhook.Message) {
	rmqQueue := conf.Config.Queue
	// 批量消息队列中
	aggregator := channelx.NewAggregator[*exhook.Message](
		func(messages []*exhook.Message) error {
			targetMessages := make([]*primitive.Message, len(messages))
			for idx, sourceMessage := range messages {
				targetMessages[idx] = buildTargetMessage(sourceMessage)
			}
			_, err := rmqProducer.SendSync(context.Background(), targetMessages...)
			if err != nil {
				log.Printf("[queue] Batch SendSync [%d] error: %s\n", len(targetMessages), err)
			}
			return nil
		},
		func(option channelx.AggregatorOption[*exhook.Message]) channelx.AggregatorOption[*exhook.Message] {
			option.BatchSize = rmqQueue.BatchSize
			option.Workers = rmqQueue.Workers
			option.ChannelBufferSize = option.BatchSize * 2
			option.LingerTime = time.Duration(rmqQueue.LingerTime) * time.Second
			log.Printf("channelx option : %v \n", option)
			return option
		},
	)
	aggregator.Start()
	defer aggregator.SafeStop()
	for {
		if sourceMessage, ok := <-ch; ok {
			aggregator.TryEnqueue(sourceMessage)
		}
	}
}

// Direct 直接发送
func Direct(rmqProducer rocketmq.Producer, ch chan *exhook.Message) {
	for {
		if sourceMessage, ok := <-ch; ok {
			targetMessages := buildTargetMessage(sourceMessage)
			err := rmqProducer.SendAsync(
				context.Background(),
				func(ctx context.Context, result *primitive.SendResult, err error) {
					if err != nil {
						log.Printf("[direct] SendAsync : %v \n", err.Error())
					}
				},
				targetMessages,
			)
			if err != nil {
				log.Printf("[direct] send message error: %s\n", err)
			}
		}
	}
}

// buildTargetMessage 构建消息
func buildTargetMessage(sourceMessage *exhook.Message) *primitive.Message {
	rmqRule := conf.Config.BridgeRule
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
	return targetMessage
}
