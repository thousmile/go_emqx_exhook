package provider

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"log"
	"strings"
	"time"
)

type RedisMessageProvider struct {
	// 目标主题
	StreamName string
	//
	PayloadFormat string
	// redis 提供者
	RedisClient redis.UniversalClient
}

func (r RedisMessageProvider) BatchSend(messages []*exhook.Message) {
	for _, message := range messages {
		r.SingleSend(message)
	}
}

func (r RedisMessageProvider) SingleSend(message *exhook.Message) {
	targetMessages := r.buildTargetMessage(message)
	err := r.RedisClient.XAdd(context.Background(), targetMessages).Err()
	if err != nil {
		log.Printf("[direct] send message error: %s\n", err)
	}
}

// BuildTargetMessage 构建消息
func (r RedisMessageProvider) buildTargetMessage(sourceMessage *exhook.Message) *redis.XAddArgs {
	values := make(map[string]interface{})
	values[SourceId] = sourceMessage.Id
	values[SourceTopic] = sourceMessage.Topic
	values[SourceNode] = sourceMessage.Node
	values[SourceFrom] = sourceMessage.From
	values[SourceQos] = int(sourceMessage.Qos)
	values[SourceTimestamp] = int64(sourceMessage.Timestamp)
	if len(sourceMessage.Headers) > 0 {
		for key, val := range sourceMessage.Headers {
			values[key] = val
		}
	}
	values[SourcePayload] = sourceMessage.Payload
	return &redis.XAddArgs{
		Stream: r.StreamName,
		Values: values,
	}
}

func BuildRedisMessageProvider(redisConf conf.RedisConfig) RedisMessageProvider {
	options := redis.UniversalOptions{
		Addrs:        redisConf.Addresses,
		DB:           redisConf.DB,
		WriteTimeout: time.Second * 1,
		ReadTimeout:  time.Second * 1,
		DialTimeout:  time.Second * 3,
		PoolSize:     4,
	}
	if len(strings.TrimSpace(redisConf.Password)) > 0 {
		options.Password = redisConf.Password
	}
	if len(strings.TrimSpace(redisConf.Username)) > 0 {
		options.Username = redisConf.Username
	}
	if len(strings.TrimSpace(redisConf.MasterName)) > 0 {
		options.MasterName = redisConf.MasterName
	}
	if len(strings.TrimSpace(redisConf.SentinelUsername)) > 0 {
		options.SentinelUsername = redisConf.SentinelUsername
	}
	if len(strings.TrimSpace(redisConf.SentinelUsername)) > 0 {
		options.SentinelUsername = redisConf.SentinelUsername
	}
	client := redis.NewUniversalClient(&options)
	return RedisMessageProvider{
		StreamName:    redisConf.StreamName,
		PayloadFormat: redisConf.PayloadFormat,
		RedisClient:   client,
	}
}
