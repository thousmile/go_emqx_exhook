package conf

import (
	"github.com/spf13/viper"
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
		},
	)
	viper.SetDefault(
		"kafkaConfig",
		KafkaConfig{
			Addresses: []string{"127.0.0.1:9092"},
			Topic:     "emqx_exhook",
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
	viper.AddConfigPath("/etc/appname/")         // 查找配置文件所在路径
	viper.AddConfigPath("$HOME/.appname")        // 多次调用AddConfigPath，可以添加多个搜索路径
	viper.AddConfigPath(".")                     // optionally look for config in the working directory
	viper.AddConfigPath("../")                   // optionally look for config in the working directory
	viper.AddConfigPath("./conf/")               // 还可以在工作目录中搜索配置文件
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Panicf("Fatal error config file: %v \n", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Panicf("Fatal error config file: %v \n", err.Error())
	}
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

	// mq 类型 Rocketmq、Rabbitmq、Kafka
	MqType string `yaml:"mqType" json:"mqType"`

	// Rocketmq 配置
	RocketmqConfig RocketmqConfig `yaml:"rocketmqConfig" json:"rocketmqConfig"`

	// Rabbitmq 配置
	RabbitmqConfig RabbitmqConfig `yaml:"rabbitmqConfig" json:"rabbitmqConfig"`

	// Kafka 配置
	KafkaConfig KafkaConfig `yaml:"kafkaConfig" json:"kafkaConfig"`

	// SendMethod 发送方式。默认是队列，queue , direct
	SendMethod string `yaml:"sendMethod" json:"sendMethod"`

	// Queue 队列配置
	Queue Queue `yaml:"queue" json:"queue"`
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
}

// KafkaConfig 桥接到 Kafka 的配置
type KafkaConfig struct {

	// Kafka Addresses
	Addresses []string `yaml:"addresses" json:"addresses"`

	// Rocketmq 的主题
	Topic string `yaml:"topic" json:"topic"`
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
