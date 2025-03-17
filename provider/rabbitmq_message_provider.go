package provider

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/wagslane/go-rabbitmq"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook_v2"
	"log"
	"os"
	"strconv"
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

func (r RabbitmqMessageProvider) BatchSend(messages []*exhook_v2.Message) {
	for _, message := range messages {
		r.SingleSend(message)
	}
}

func (r RabbitmqMessageProvider) SingleSend(message *exhook_v2.Message) {
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
func (r RabbitmqMessageProvider) buildTargetMessageHeaders(sourceMessage *exhook_v2.Message) rabbitmq.Table {
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
	return headers
}

func BuildRabbitmqMessageProvider(rbbConf conf.RabbitmqConfig) RabbitmqMessageProvider {
	resolver := rabbitmq.NewStaticResolver(rbbConf.Addresses, false)
	tlsConf := createRabbitTLS(rbbConf.Tls)
	conn, err := rabbitmq.NewClusterConn(resolver,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsConfig(
			rabbitmq.Config{TLSClientConfig: tlsConf},
		),
	)
	if err != nil {
		log.Panicf("rabbitmq conn error %v", err)
	}
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(rbbConf.ExchangeName),
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)
	if err != nil {
		log.Panicf("rabbitmq producer error %v", err)
	}
	p1 := RabbitmqMessageProvider{
		RoutingKeys:    rbbConf.RoutingKeys,
		ExchangeName:   rbbConf.ExchangeName,
		RabbitProducer: publisher,
	}
	return p1
}

func createRabbitTLS(sasl conf.TlsConfig) (t *tls.Config) {
	if sasl.CertFile != "" && sasl.KeyFile != "" && sasl.CaFile != "" {
		cert, err := tls.LoadX509KeyPair(sasl.CertFile, sasl.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		caCert, err := os.ReadFile(sasl.CaFile)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: sasl.TlsSkipVerify,
		}
	}
	return t
}
