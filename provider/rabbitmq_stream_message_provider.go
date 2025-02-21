package provider

import (
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/message"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"log"
	"strconv"
	"time"
)

var rabbitmqStreamCodecs = map[string]stream.Compression{
	"none":   stream.Compression{}.None(),
	"gzip":   stream.Compression{}.Gzip(),
	"snappy": stream.Compression{}.Snappy(),
	"lz4":    stream.Compression{}.Lz4(),
	"zstd":   stream.Compression{}.Zstd(),
}

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
	maxAge := 7 * 24 * time.Hour
	if len(rbbConf.MaxAge) > 0 {
		maxAge, err = time.ParseDuration(rbbConf.MaxAge)
		if err != nil {
			log.Panicf("rabbitmq stream maxAge format error %v", err)
		}
	}
	maxLengthBytes := stream.ByteCapacity{}.GB(10)
	if len(rbbConf.MaxLengthBytes) > 0 {
		maxLengthBytes = stream.ByteCapacity{}.From(rbbConf.MaxLengthBytes)
	}
	maxSegmentSizeBytes := stream.ByteCapacity{}.GB(1)
	if len(rbbConf.MaxSegmentSizeBytes) > 0 {
		maxSegmentSizeBytes = stream.ByteCapacity{}.From(rbbConf.MaxSegmentSizeBytes)
	}
	err = env.DeclareStream(rbbConf.StreamName,
		stream.NewStreamOptions().
			SetMaxAge(maxAge).
			SetMaxLengthBytes(maxLengthBytes).
			SetMaxSegmentSizeBytes(maxSegmentSizeBytes),
	)
	defCodec := stream.Compression{}.None()
	if len(rbbConf.CompressionCodec) > 0 {
		codec1, ok := rabbitmqStreamCodecs[rbbConf.CompressionCodec]
		if ok {
			defCodec = codec1
		}
	}
	producer, err := env.NewProducer(rbbConf.StreamName,
		stream.NewProducerOptions().
			SetCompression(defCodec).
			SetSubEntrySize(100),
	)
	if err != nil {
		log.Panicf("rabbitmq stream producer error %v", err)
	}
	p1 := RabbitmqStreamMessageProvider{
		RabbitStreamProducer: producer,
	}
	return p1
}
