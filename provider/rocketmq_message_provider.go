package provider

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook_v2"
	"log"
	"strconv"
	"time"
)

type RocketmqMessageProvider struct {
	// 目标主题
	Topic string
	// 目标tag
	Tag string
	// rocketmq 提供者
	RmqProducer rocketmq.Producer
}

func (r RocketmqMessageProvider) BatchSend(messages []*exhook_v2.Message) {
	targetMessages := make([]*primitive.Message, len(messages))
	for idx, sourceMessage := range messages {
		targetMessages[idx] = r.buildTargetMessage(sourceMessage)
	}
	_, err := r.RmqProducer.SendSync(context.Background(), targetMessages...)
	if err != nil {
		log.Printf("[queue] rocketmq batch send [%d] error: %s\n", len(targetMessages), err)
	}
}

func (r RocketmqMessageProvider) SingleSend(message *exhook_v2.Message) {
	targetMessages := r.buildTargetMessage(message)
	err := r.RmqProducer.SendAsync(
		context.Background(),
		func(ctx context.Context, result *primitive.SendResult, err error) {
			if err != nil {
				log.Printf("[direct] rocketmq single send error: %v \n", err.Error())
			}
		},
		targetMessages,
	)
	if err != nil {
		log.Printf("[direct] rocketmq send message error: %s\n", err)
	}
}

// BuildTargetMessage 构建消息
func (r RocketmqMessageProvider) buildTargetMessage(sourceMessage *exhook_v2.Message) *primitive.Message {
	targetMessage := &primitive.Message{
		Topic: r.Topic,
		Body:  sourceMessage.Payload,
	}
	targetMessage.WithKeys([]string{sourceMessage.Id})
	targetMessage.WithTag(r.Tag)
	targetMessage.WithProperty(SourceId, sourceMessage.Id)
	targetMessage.WithProperty(SourceTopic, sourceMessage.Topic)
	targetMessage.WithProperty(SourceNode, sourceMessage.Node)
	targetMessage.WithProperty(SourceFrom, sourceMessage.From)
	targetMessage.WithProperty(SourceQos, strconv.Itoa(int(sourceMessage.Qos)))
	targetMessage.WithProperty(SourceTimestamp,
		strconv.FormatInt(int64(sourceMessage.Timestamp), 10))
	if len(sourceMessage.Headers) > 0 {
		for key, val := range sourceMessage.Headers {
			targetMessage.WithProperty(key, val)
		}
	}
	return targetMessage
}

func BuildRocketmqMessageProvider(rmqConf conf.RocketmqConfig) RocketmqMessageProvider {
	var acl primitive.Credentials
	if len(rmqConf.AccessKey) > 0 && len(rmqConf.SecretKey) > 0 {
		acl.AccessKey = rmqConf.AccessKey
		acl.SecretKey = rmqConf.SecretKey
	}
	rmqProducer, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(rmqConf.NameServer)),
		producer.WithGroupName(rmqConf.GroupName),
		producer.WithSendMsgTimeout(time.Second*1),
		producer.WithCredentials(acl),
	)
	err := rmqProducer.Start()
	if err != nil {
		log.Panicf("rocketmq producer error %v", err)
	}
	p1 := RocketmqMessageProvider{
		Topic:       rmqConf.Topic,
		Tag:         rmqConf.Tag,
		RmqProducer: rmqProducer,
	}
	return p1
}
