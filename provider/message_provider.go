package provider

import "go_emqx_exhook/emqx.io/grpc/exhook"

const (
	SourceId        = "sourceId"
	SourceTopic     = "sourceTopic"
	SourceNode      = "sourceNode"
	SourceFrom      = "sourceFrom"
	SourceQos       = "sourceQos"
	SourceTimestamp = "sourceTimestamp"
)

// MessageProvider 抽象出来的 消息提供者
type MessageProvider interface {

	// BatchSend 批量发送
	BatchSend(messages []*exhook.Message)

	// SingleSend 单条发送
	SingleSend(message *exhook.Message)
}
