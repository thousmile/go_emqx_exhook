package conf

import (
	"errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
)

// Config 全局配置配置文件
var Config *ServerConfig

func init() {
	viper.SetDefault("appName", "go_emqx_exhook")
	viper.SetDefault("port", 16565)
	viper.SetDefault("chanBufferSize", 10240)
	viper.SetDefault("mqType", "Rocketmq")
	viper.SetDefault(
		"bridgeRule",
		BridgeRule{
			Topics: []string{"/#"},
		},
	)
	viper.SetDefault(
		"rocketmqConfig",
		RocketmqConfig{
			NameServer: []string{
				"127.0.0.1:9876",
			},
			Topic:     "emqx_exhook",
			Tag:       "exhook",
			GroupName: "exhook",
		},
	)
	viper.SetDefault(
		"rabbitmqConfig",
		RabbitmqConfig{
			Addresses:    []string{"amqp://guest:guest@127.0.0.1:5672"},
			ExchangeName: "amq.direct",
			RoutingKeys:  []string{"exhook"},
			Tls: TlsConfig{
				Enable:        false,
				TlsSkipVerify: false,
				CaFile:        "certs/ca/ca.crt",
				CertFile:      "certs/client/client.crt",
				KeyFile:       "certs/client/client.key",
			},
		},
	)
	viper.SetDefault(
		"kafkaConfig",
		KafkaConfig{
			Addresses: []string{"127.0.0.1:9092"},
			Topic:     "emqx_exhook",
			Sasl: KafkaSasl{
				Enable:    false,
				User:      "exhook",
				Password:  "exhook",
				Algorithm: "plain",
			},
			Tls: TlsConfig{
				Enable:        false,
				TlsSkipVerify: false,
				CaFile:        "certs/ca/ca.crt",
				CertFile:      "certs/client/client.crt",
				KeyFile:       "certs/client/client.key",
			},
		},
	)
	viper.SetDefault(
		"redisConfig",
		RedisConfig{
			Addresses:  []string{"127.0.0.1:6379"},
			StreamName: "emqx_exhook",
			DB:         0,
		},
	)
	viper.SetDefault("sendMethod", "queue")
	viper.SetDefault(
		"queue",
		Queue{
			BatchSize:  100,
			Workers:    2,
			LingerTime: 1,
		},
	)
	viper.SetConfigName("config")                // name of config file (without extension)
	viper.SetConfigType("yaml")                  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                     // 程序所在路径
	viper.AddConfigPath("/etc/go_emqx_exhook/")  // Linux /etc/go_emqx_exhook 目录下
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Panicf("viper ReadInConfig error : %v \n", err)
		}
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Panicf("viper Unmarshal error : %v \n", err.Error())
	}
	out, _ := yaml.Marshal(Config)
	log.Printf("app config : \n%s\n", string(out))
}

type ServerConfig struct {
	// 服务名称 ，默认: go-emqx-exhook
	AppName string `yaml:"appName" json:"appName"`

	// 服务端口 ，默认: 16565
	Port int `yaml:"port" json:"port"`

	// ChanBufferSize 管道缓冲区大小，默认: 10240
	ChanBufferSize int `yaml:"chanBufferSize" json:"chanBufferSize"`

	// 桥接规则
	BridgeRule BridgeRule `yaml:"bridgeRule" json:"bridgeRule"`

	// mq 类型 Rocketmq、Rabbitmq、Kafka、Redis
	MqType string `yaml:"mqType" json:"mqType"`

	// Rocketmq 配置
	RocketmqConfig RocketmqConfig `yaml:"rocketmqConfig" json:"rocketmqConfig"`

	// Rabbitmq 配置
	RabbitmqConfig RabbitmqConfig `yaml:"rabbitmqConfig" json:"rabbitmqConfig"`

	// Kafka 配置
	KafkaConfig KafkaConfig `yaml:"kafkaConfig" json:"kafkaConfig"`

	// Redis 配置
	RedisConfig RedisConfig `yaml:"redisConfig" json:"redisConfig"`

	// SendMethod 发送方式。默认是队列，queue , direct
	SendMethod string `yaml:"sendMethod" json:"sendMethod"`

	// Queue 队列配置
	Queue Queue `yaml:"queue" json:"queue"`

	// grpc server tls
	Tls TlsConfig `yaml:"tls" json:"tls"`
}

// BridgeRule 桥接规则
type BridgeRule struct {
	// Emqx 的主题
	Topics []string `yaml:"topics" json:"topics"`
}

// RocketmqConfig 桥接到 Rocketmq 的配置
type RocketmqConfig struct {

	// Rocketmq NameServer
	NameServer []string `yaml:"nameServer" json:"nameServer"`

	// Rocketmq 的主题
	Topic string `yaml:"topic" json:"topic"`

	// Rocketmq 的 Tag
	Tag string `yaml:"tag" json:"tag"`

	// Rocketmq 的 分组名称
	GroupName string `yaml:"groupName" json:"groupName"`

	// 阿里云 Rocketmq 的 accessKey
	AccessKey string `yaml:"accessKey" json:"accessKey"`

	// 阿里云 Rocketmq 的 secretKey
	SecretKey string `yaml:"secretKey" json:"secretKey"`
}

// RabbitmqConfig 桥接到 Rabbitmq 的配置
type RabbitmqConfig struct {

	// Rocketmq Addresses
	Addresses []string `yaml:"addresses" json:"addresses"`

	// Rocketmq RoutingKeys
	RoutingKeys []string `yaml:"routingKeys" json:"routingKeys"`

	// Rocketmq ExchangeName
	ExchangeName string `yaml:"exchangeName" json:"exchangeName"`

	// Rocketmq tls
	Tls TlsConfig `yaml:"tls" json:"tls"`
}

// KafkaConfig 桥接到 Kafka 的配置
type KafkaConfig struct {

	// Kafka Addresses
	Addresses []string `yaml:"addresses" json:"addresses"`

	// Kafka 的主题
	Topic string `yaml:"topic" json:"topic"`

	// Kafka SASL
	Sasl KafkaSasl `yaml:"sasl" json:"sasl"`

	// Kafka tls
	Tls TlsConfig `yaml:"tls" json:"tls"`
}

// KafkaSasl 配置
type KafkaSasl struct {
	// 启用
	Enable bool `yaml:"enable" json:"enable"`

	// 用户名
	User string `yaml:"user" json:"user"`

	// 密码
	Password string `yaml:"password" json:"password"`

	// The SASL SCRAM SHA algorithm sha256 or sha512 as mechanism
	Algorithm string `yaml:"algorithm" json:"algorithm"`
}

// TlsConfig 配置
type TlsConfig struct {
	// 启用
	Enable bool `yaml:"enable" json:"enable"`

	// Whether to skip TLS server cert verification
	TlsSkipVerify bool `yaml:"tlsSkipVerify" json:"tlsSkipVerify"`

	// The optional certificate authority file for TLS client authentication
	CaFile string `yaml:"caFile" json:"caFile"`

	// The optional certificate file for client authentication
	CertFile string `yaml:"certFile" json:"certFile"`

	// The optional key file for client authentication
	KeyFile string `yaml:"keyFile" json:"keyFile"`
}

type RedisConfig struct {
	// redis 地址。默认: 127.0.0.1:6379
	Addresses []string `yaml:"addresses" json:"addresses"`

	// redis stream
	StreamName string `yaml:"streamName" json:"streamName"`

	// 用户名，默认: 空
	Username string `yaml:"username" json:"username"`

	// 密码，默认: 空
	Password string `yaml:"password" json:"password"`

	// 库索引，默认: 0
	DB int `yaml:"db" json:"db"`

	// Sentinel 模式。
	MasterName string `yaml:"masterName" json:"masterName"`

	// Sentinel username。
	SentinelUsername string `yaml:"sentinelUsername" json:"sentinelUsername"`

	// Sentinel password
	SentinelPassword string `yaml:"sentinelPassword" json:"sentinelPassword"`
}

// Queue 配置
type Queue struct {
	// 批量处理的数据，默认: 100 条
	BatchSize int `yaml:"batchSize" json:"batchSize"`

	// 工作线程，默认: 2 个
	Workers int `yaml:"workers" json:"workers"`

	// LingerTime 延迟时间(单位秒)，默认: 1 秒
	LingerTime int `yaml:"lingerTime" json:"lingerTime"`
}
