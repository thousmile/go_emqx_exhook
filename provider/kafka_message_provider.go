package provider

import (
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/xdg-go/scram"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"log"
	"os"
	"strconv"
	"time"
)

var kafkaCodecs = map[string]sarama.CompressionCodec{
	"none":   sarama.CompressionNone,
	"gzip":   sarama.CompressionGZIP,
	"snappy": sarama.CompressionSnappy,
	"lz4":    sarama.CompressionLZ4,
	"zstd":   sarama.CompressionZSTD,
}

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
		log.Printf("[queue] kafka batch send [%d] error: %s\n", len(targetMessages), err)
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
	timestamp := int64(sourceMessage.Timestamp)
	headers := []sarama.RecordHeader{
		{Key: []byte(SourceId), Value: []byte(sourceMessage.Id)},
		{Key: []byte(SourceTopic), Value: []byte(sourceMessage.Topic)},
		{Key: []byte(SourceNode), Value: []byte(sourceMessage.Node)},
		{Key: []byte(SourceFrom), Value: []byte(sourceMessage.From)},
		{Key: []byte(SourceQos), Value: []byte(strconv.Itoa(int(sourceMessage.Qos)))},
		{Key: []byte(SourceTimestamp), Value: []byte(strconv.FormatInt(timestamp, 10))},
	}
	if len(sourceMessage.Headers) > 0 {
		for key, val := range sourceMessage.Headers {
			header := sarama.RecordHeader{Key: []byte(key), Value: []byte(val)}
			headers = append(headers, header)
		}
	}
	return &sarama.ProducerMessage{
		Topic:     r.Topic,
		Key:       sarama.StringEncoder(sourceMessage.Id),
		Value:     sarama.ByteEncoder(sourceMessage.Payload),
		Headers:   headers,
		Timestamp: time.UnixMilli(timestamp),
	}
}

func BuildKafkaMessageProvider(kafkaConf conf.KafkaConfig) KafkaMessageProvider {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //ACK,发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //分区,新选出一个分区
	config.Producer.Return.Successes = true                   //确认,成功交付的消息将在success channel返回
	codec, ok := kafkaCodecs[kafkaConf.CompressionCodec]
	if !ok {
		codec = sarama.CompressionNone
	}
	config.Producer.Compression = codec
	if kafkaConf.Sasl.Enable {
		sasl := kafkaConf.Sasl
		config.Net.SASL.Enable = true
		config.Net.SASL.User = sasl.User
		config.Net.SASL.Password = sasl.Password
		config.Net.SASL.Handshake = true
		switch sasl.Algorithm {
		case "sha512":
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		case "sha256":
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		default:
			config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		}
	}
	if kafkaConf.Tls.Enable {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = createKafkaTLS(kafkaConf.Tls)
	}
	client, err := sarama.NewClient(kafkaConf.Addresses, config)
	if err != nil {
		log.Panicf("kafka client error %v", err)
	}
	clusterAdmin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		log.Panicf("kafka cluster admin error %v", err)
	}
	topics, err := clusterAdmin.ListTopics()
	if err != nil {
		log.Panicf("kafka list topics error %v", err)
	}
	_, ok = topics[kafkaConf.Topic]
	if !ok {
		// 主题不存在，就创建主题
		detail := &sarama.TopicDetail{NumPartitions: -1, ReplicationFactor: -1}
		if kafkaConf.NumPartitions > -1 {
			detail.NumPartitions = kafkaConf.NumPartitions
		}
		if kafkaConf.ReplicationFactor > -1 {
			detail.ReplicationFactor = kafkaConf.ReplicationFactor
		}
		if len(kafkaConf.ConfigEntries) > 0 {
			entries := make(map[string]*string, len(kafkaConf.ConfigEntries))
			for k, v := range kafkaConf.ConfigEntries {
				str := fmt.Sprintf("%v", v)
				entries[k] = &str
			}
			detail.ConfigEntries = entries
		}
		err = clusterAdmin.CreateTopic(kafkaConf.Topic, detail, false)
		if err != nil {
			log.Panicf("kafka create topic %s error %v", kafkaConf.Topic, err)
		}
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Panicf("kafka producer error %v", err)
	}
	p1 := KafkaMessageProvider{
		Topic:         kafkaConf.Topic,
		KafkaProducer: producer,
	}
	return p1
}

func createKafkaTLS(sasl conf.TlsConfig) (t *tls.Config) {
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

var (
	SHA256 scram.HashGeneratorFcn = sha256.New
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
