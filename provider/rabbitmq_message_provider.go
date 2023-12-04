package provider

import (
	"github.com/wagslane/go-rabbitmq"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"log"
	"strings"
)

type RabbitmqMessageProvider struct {
	// 目标主题
	TargetTopic string
	// rabbitmq 提供者
	RabbitProducer *rabbitmq.Publisher
}

func (r RabbitmqMessageProvider) BatchSend(messages []*exhook.Message) {
	//TODO implement me
	panic("implement me")
}

func (r RabbitmqMessageProvider) SingleSend(message *exhook.Message) {

	//TODO implement me
	panic("implement me")
}

// BuildTargetMessage 构建消息
func (r RabbitmqMessageProvider) buildTargetMessage(sourceMessage *exhook.Message) []byte {
	return nil
}

func BuildRabbitmqMessageProvider(rbbConf conf.RabbitmqConfig, targetTopic string) RabbitmqMessageProvider {
	url := strings.Join(rbbConf.Addresses, ",")
	conn, err := rabbitmq.NewConn(
		url,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *rabbitmq.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(rbbConf.ExchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	defer publisher.Close()

	publisher.NotifyReturn(func(r rabbitmq.Return) {
		log.Printf("message returned from server: %s", string(r.Body))
	})

	publisher.NotifyPublish(func(c rabbitmq.Confirmation) {
		log.Printf("message confirmed from server. tag: %v, ack: %v", c.DeliveryTag, c.Ack)
	})

	if err != nil {
		log.Fatal(err)
	}
	p1 := RabbitmqMessageProvider{
		TargetTopic:    targetTopic,
		RabbitProducer: publisher,
	}
	return p1
}
