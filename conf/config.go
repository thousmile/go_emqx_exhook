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
	viper.SetDefault(
		"bridgeRule",
		BridgeRule{
			SourceTopics: []string{
				"/#",
			},
			TargetTopic: "emqx_msg_bridge",
			TargetTag:   "emqx",
		},
	)
	viper.SetDefault(
		"rocketmqConfig",
		RocketmqConfig{
			NameServer: []string{
				"127.0.0.1:9876",
			},
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
		log.Panicf("Fatal error config file: %v \n", err)
	}
}

type ServerConfig struct {
	// 服务名称
	AppName string `yaml:"appName" json:"appName"`

	// 服务端口
	Port int `yaml:"port" json:"port"`

	// 桥接规则
	BridgeRule BridgeRule `yaml:"bridgeRule" json:"bridgeRule"`

	// Rocketmq 配置
	RocketmqConfig RocketmqConfig `yaml:"rocketmqConfig" json:"rocketmqConfig"`
}

// RocketmqConfig 桥接到 Rocketmq 的配置
type RocketmqConfig struct {

	// Rocketmq NameServer
	NameServer []string `yaml:"nameServer" json:"nameServer"`
}

// BridgeRule 桥接规则
type BridgeRule struct {
	// Emqx 的主题
	SourceTopics []string `yaml:"sourceTopics" json:"sourceTopics"`

	// Rocketmq 的主题
	TargetTopic string `yaml:"targetTopic" json:"targetTopic"`

	// Rocketmq 的 Tag
	TargetTag string `yaml:"targetTag" json:"targetTag"`
}
