package provider

import (
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/message"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"log"
	"strconv"
)

type RabbitmqStreamMessageProvider struct {
	// rabbitmq stream 提供者
	RabbitStreamProducer *stream.Producer
}

func (r RabbitmqStreamMessageProvider) BatchSend(messages []*exhook.Message) {
	targetMessages := make([]message.StreamMessage, len(messages))
	for idx, sourceMessage := range messages {
		targetMessages[idx] = r.buildTargetMessage(sourceMessage)
	}
	err := r.RabbitStreamProducer.BatchSend(targetMessages)
	if err != nil {
		log.Printf("[queue] rabbitmq stream batch send [%d] error: %s\n", len(targetMessages), err)
	}
}

func (r RabbitmqStreamMessageProvider) SingleSend(message *exhook.Message) {
	targetMessage := r.buildTargetMessage(message)
	err := r.RabbitStreamProducer.Send(targetMessage)
	if err != nil {
		log.Printf("[direct] rabbitmq stream single send error: %v \n", err.Error())
	}
}

// BuildTargetMessage 构建消息
func (r RabbitmqStreamMessageProvider) buildTargetMessage(sourceMessage *exhook.Message) message.StreamMessage {
	streamMessage := amqp.NewMessage(sourceMessage.GetPayload())
	headers := map[string]interface{}{
		SourceId:        sourceMessage.Id,
		SourceTopic:     sourceMessage.Topic,
		SourceNode:      sourceMessage.Node,
		SourceFrom:      sourceMessage.From,
		SourceQos:       strconv.Itoa(int(sourceMessage.Qos)),
		SourceTimestamp: strconv.FormatInt(int64(sourceMessage.Timestamp), 10),
	}
	if len(sourceMessage.Headers) > 0 {
		for key, val := range sourceMessage.Headers {
			headers[key] = val
		}
	}
	streamMessage.ApplicationProperties = headers
	return streamMessage
}

func BuildRabbitmqStreamMessageProvider(rbbConf conf.RabbitmqStreamConfig) RabbitmqStreamMessageProvider {
	tlsConf := createRabbitTLS(rbbConf.Tls)
	if rbbConf.MaxProducersPerClient < 1 {
		rbbConf.MaxProducersPerClient = 2
	}
	options := stream.NewEnvironmentOptions().
		SetUris(rbbConf.Addresses).
		SetMaxProducersPerClient(rbbConf.MaxProducersPerClient).
		SetTLSConfig(tlsConf)
	env, err := stream.NewEnvironment(options)
	producer, err := env.NewProducer(rbbConf.StreamName, stream.NewProducerOptions())
	if err != nil {
		log.Panicf("rabbitmq stream producer error %v", err)
	}
	p1 := RabbitmqStreamMessageProvider{
		RabbitStreamProducer: producer,
	}
	return p1
}
