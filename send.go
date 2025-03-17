package main

import (
	"go_emqx_exhook/channelx"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook_v2"
	"go_emqx_exhook/provider"
	"log"
	"time"
)

// Queue 使用队列
func Queue(producer provider.MessageProvider, ch chan *exhook_v2.Message) {
	queue := conf.Config.Queue
	// 批量消息队列中
	aggregator := channelx.NewAggregator[*exhook_v2.Message](
		func(messages []*exhook_v2.Message) error {
			producer.BatchSend(messages)
			return nil
		},
		func(option channelx.AggregatorOption[*exhook_v2.Message]) channelx.AggregatorOption[*exhook_v2.Message] {
			option.BatchSize = queue.BatchSize
			option.Workers = queue.Workers
			option.ChannelBufferSize = option.BatchSize * 2
			option.LingerTime = time.Duration(queue.LingerTime) * time.Second
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
func Direct(producer provider.MessageProvider, ch chan *exhook_v2.Message) {
	for {
		if sourceMessage, ok := <-ch; ok {
			producer.SingleSend(sourceMessage)
		}
	}
}
