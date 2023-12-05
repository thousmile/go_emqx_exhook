package provider

import (
	"github.com/wagslane/go-rabbitmq"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"log"
	"strconv"
	"strings"
	"time"
)

type RabbitmqMessageProvider struct {
	// 目标主题
	RoutingKeys []string
	// 交换机
	ExchangeName string
	// rabbitmq 提供者
	RabbitProducer *rabbitmq.Publisher
	// rabbitmq 连接
	RabbitmqConn *rabbitmq.Conn
}

func (r RabbitmqMessageProvider) BatchSend(messages []*exhook.Message) {
	for _, message := range messages {
		r.SingleSend(message)
	}
}

func (r RabbitmqMessageProvider) SingleSend(message *exhook.Message) {
	headers := r.buildTargetMessageHeaders(message)
	err := r.RabbitProducer.Publish(
		message.Payload,
		r.RoutingKeys,
		rabbitmq.WithPublishOptionsMessageID(message.Id),
		rabbitmq.WithPublishOptionsTimestamp(time.UnixMilli(int64(message.Timestamp))),
		rabbitmq.WithPublishOptionsHeaders(headers),
		rabbitmq.WithPublishOptionsExchange(r.ExchangeName),
	)
	if err != nil {
		log.Printf("[direct] rabbitmq single send error: %v \n", err.Error())
	}
}

// BuildTargetMessage 构建消息
func (r RabbitmqMessageProvider) buildTargetMessageHeaders(sourceMessage *exhook.Message) rabbitmq.Table {
	headers := map[string]interface{}{
		SourceId:        sourceMessage.Id,
		SourceTopic:     sourceMessage.Topic,
		SourceNode:      sourceMessage.Node,
		SourceFrom:      sourceMessage.From,
		SourceQos:       strconv.Itoa(int(sourceMessage.Qos)),
		SourceTimestamp: strconv.FormatInt(int64(sourceMessage.Timestamp), 10),
	}
	return headers
}

func BuildRabbitmqMessageProvider(rbbConf conf.RabbitmqConfig) RabbitmqMessageProvider {
	url := strings.Join(rbbConf.Addresses, ",")
	conn, err := rabbitmq.NewConn(url,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(rbbConf.ExchangeName),
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)
	if err != nil {
		log.Fatal(err)
	}
	p1 := RabbitmqMessageProvider{
		RoutingKeys:    rbbConf.RoutingKeys,
		ExchangeName:   rbbConf.ExchangeName,
		RabbitProducer: publisher,
	}
	return p1
}
