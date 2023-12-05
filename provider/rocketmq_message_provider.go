package provider

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
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

func (r RocketmqMessageProvider) BatchSend(messages []*exhook.Message) {
	targetMessages := make([]*primitive.Message, len(messages))
	for idx, sourceMessage := range messages {
		targetMessages[idx] = r.buildTargetMessage(sourceMessage)
	}
	_, err := r.RmqProducer.SendSync(context.Background(), targetMessages...)
	if err != nil {
		log.Printf("[queue] rocketmq batch send [%d] error: %s\n", len(targetMessages), err)
	}
}

func (r RocketmqMessageProvider) SingleSend(message *exhook.Message) {
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
		log.Printf("[direct] send message error: %s\n", err)
	}
}

// BuildTargetMessage 构建消息
func (r RocketmqMessageProvider) buildTargetMessage(sourceMessage *exhook.Message) *primitive.Message {
	targetMessage := &primitive.Message{
		Topic: r.Topic,
		Body:  sourceMessage.Payload,
	}
	targetMessage.WithKeys([]string{sourceMessage.Id})
	targetMessage.WithTag(r.Tag)
	for key, val := range sourceMessage.GetHeaders() {
		targetMessage.WithProperty(key, val)
	}
	targetMessage.WithProperty(SourceId, sourceMessage.Id)
	targetMessage.WithProperty(SourceTopic, sourceMessage.Topic)
	targetMessage.WithProperty(SourceNode, sourceMessage.Node)
	targetMessage.WithProperty(SourceFrom, sourceMessage.From)
	targetMessage.WithProperty(SourceQos, strconv.Itoa(int(sourceMessage.Qos)))
	targetMessage.WithProperty(SourceTimestamp,
		strconv.FormatInt(int64(sourceMessage.Timestamp), 10))
	return targetMessage
}

func BuildRocketmqMessageProvider(rmqConf conf.RocketmqConfig) RocketmqMessageProvider {
	rmqProducer, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(rmqConf.NameServer)),
		producer.WithGroupName(rmqConf.GroupName),
		producer.WithSendMsgTimeout(time.Second*1),
	)
	err := rmqProducer.Start()
	if err != nil {
		log.Fatal(err)
	}
	p1 := RocketmqMessageProvider{
		Topic:       rmqConf.Topic,
		Tag:         rmqConf.Tag,
		RmqProducer: rmqProducer,
	}
	return p1
}
