package provider

import (
	"fmt"
	"github.com/IBM/sarama"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"log"
	"strconv"
	"time"
)

type KafkaMessageProvider struct {
	// 目标主题
	Topic string
	// kafka 提供者
	KafkaProducer sarama.SyncProducer
}

func (r KafkaMessageProvider) BatchSend(messages []*exhook.Message) {
	targetMessages := make([]*sarama.ProducerMessage, len(messages))
	for idx, sourceMessage := range messages {
		targetMessages[idx] = r.buildTargetMessage(sourceMessage)
	}
	err := r.KafkaProducer.SendMessages(targetMessages)
	if err != nil {
		log.Printf("[queue] rocketmq batch send [%d] error: %s\n", len(targetMessages), err)
	}
}

func (r KafkaMessageProvider) SingleSend(message *exhook.Message) {
	targetMessage := r.buildTargetMessage(message)
	_, _, err := r.KafkaProducer.SendMessage(targetMessage)
	if err != nil {
		log.Printf("[direct] kafka single send error: %v \n", err.Error())
	}
}

// BuildTargetMessage 构建消息
func (r KafkaMessageProvider) buildTargetMessage(sourceMessage *exhook.Message) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: r.Topic,
		Key:   sarama.StringEncoder(sourceMessage.Id),
		Value: sarama.ByteEncoder(sourceMessage.Payload),
		Headers: []sarama.RecordHeader{
			{Key: []byte(SourceId), Value: []byte(sourceMessage.Id)},
			{Key: []byte(SourceTopic), Value: []byte(sourceMessage.Topic)},
			{Key: []byte(SourceNode), Value: []byte(sourceMessage.Node)},
			{Key: []byte(SourceFrom), Value: []byte(sourceMessage.From)},
			{Key: []byte(SourceQos), Value: []byte(strconv.Itoa(int(sourceMessage.Qos)))},
			{Key: []byte(SourceTimestamp), Value: []byte(strconv.FormatInt(int64(sourceMessage.Timestamp), 10))},
		},
		Timestamp: time.UnixMilli(int64(sourceMessage.Timestamp)),
	}
}

func BuildKafkaMessageProvider(kafkaConf conf.KafkaConfig) KafkaMessageProvider {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //ACK,发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //分区,新选出一个分区
	config.Producer.Return.Successes = true                   //确认,成功交付的消息将在success channel返回
	client, err := sarama.NewSyncProducer(kafkaConf.Addresses, config)
	if err != nil {
		fmt.Println("Producer error", err)
	}
	p1 := KafkaMessageProvider{
		Topic:         kafkaConf.Topic,
		KafkaProducer: client,
	}
	return p1
}
