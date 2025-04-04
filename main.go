package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"go_emqx_exhook/conf"
	"go_emqx_exhook/emqx.io/grpc/exhook_v2"
	"go_emqx_exhook/emqx.io/grpc/exhook_v3"
	"go_emqx_exhook/impl"
	"go_emqx_exhook/provider"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	appConf := conf.Config
	rule := appConf.BridgeRule
	// 创建一个消息提供者
	var msgProvider provider.MessageProvider
	// 关闭连接
	connClose := func(conn io.Closer) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	if strings.EqualFold(appConf.MqType, "Rabbitmq") || strings.EqualFold(appConf.MqType, "rabbitmq") {
		rabbit := provider.BuildRabbitmqMessageProvider(appConf.RabbitmqConfig)
		defer rabbit.RabbitProducer.Close()
		defer connClose(rabbit.RabbitmqConn)
		msgProvider = rabbit
	} else if strings.EqualFold(appConf.MqType, "RabbitmqStream") || strings.EqualFold(appConf.MqType, "rabbitmqStream") {
		rabbitMQStream := provider.BuildRabbitmqStreamMessageProvider(appConf.RabbitmqStreamConfig)
		defer connClose(rabbitMQStream.RabbitStreamProducer)
		msgProvider = rabbitMQStream
	} else if strings.EqualFold(appConf.MqType, "Kafka") || strings.EqualFold(appConf.MqType, "kafka") {
		kafka := provider.BuildKafkaMessageProvider(appConf.KafkaConfig)
		defer connClose(kafka.KafkaProducer)
		msgProvider = kafka
	} else if strings.EqualFold(appConf.MqType, "Redis") || strings.EqualFold(appConf.MqType, "redis") {
		redisMq := provider.BuildRedisMessageProvider(appConf.RedisConfig)
		defer connClose(redisMq.RedisClient)
		msgProvider = redisMq
	} else {
		rmq := provider.BuildRocketmqMessageProvider(appConf.RocketmqConfig)
		defer func(p rocketmq.Producer) {
			err := p.Shutdown()
			if err != nil {
				log.Fatal(err)
			}
		}(rmq.RmqProducer)
		msgProvider = rmq
	}

	ch := make(chan *exhook_v2.Message, appConf.ChanBufferSize)

	// 发送方式“ queue or direct ”
	if appConf.SendMethod == "queue" {
		go Queue(msgProvider, ch)
	} else {
		go Direct(msgProvider, ch)
	}

	var grpcServerOptions []grpc.ServerOption
	tlsCfg := appConf.Tls
	if tlsCfg.Enable {
		grpcServerOptions = append(grpcServerOptions, grpc.Creds(getServerCred(tlsCfg)))
	}
	srv := grpc.NewServer(grpcServerOptions...)

	// 注册 emqx 的 exhook v2 grpc 服务
	exhook_v2.RegisterHookProviderServer(srv, &impl.HookProviderServerV2Impl{
		SourceTopics: rule.Topics,
		Callback: func(request *exhook_v2.MessagePublishRequest) {
			ch <- request.GetMessage()
		},
	})

	// 注册 emqx 的 exhook v3 grpc 服务
	exhook_v3.RegisterHookProviderServer(srv, &impl.HookProviderServerV3Impl{
		SourceTopics: rule.Topics,
		Callback: func(request *exhook_v3.MessagePublishRequest) {
			ch <- &exhook_v2.Message{
				Node:      request.GetMessage().GetNode(),
				Id:        request.GetMessage().GetId(),
				Qos:       request.GetMessage().GetQos(),
				From:      request.GetMessage().GetFrom(),
				Topic:     request.GetMessage().GetTopic(),
				Payload:   request.GetMessage().GetPayload(),
				Timestamp: request.GetMessage().GetTimestamp(),
				Headers:   request.GetMessage().GetHeaders(),
			}
		},
	})

	// 监听 指定端口的 tcp 连接
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", appConf.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(lis net.Listener) {
		err := lis.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(lis)
	log.Printf("%s [%s] %s => grpc server listen port : %d \n", appConf.AppName, appConf.SendMethod, appConf.MqType, appConf.Port)
	_ = srv.Serve(lis)
}

func getServerCred(tlsCfg conf.TlsConfig) credentials.TransportCredentials {
	cert, _ := tls.LoadX509KeyPair(tlsCfg.CertFile, tlsCfg.KeyFile)
	certPool := x509.NewCertPool()
	ca, _ := os.ReadFile(tlsCfg.CaFile)
	certPool.AppendCertsFromPEM(ca)
	cred := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	return cred
}
