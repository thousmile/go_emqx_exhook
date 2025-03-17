package provider

import (
	"go_emqx_exhook/emqx.io/grpc/exhook_v2"
)

const (
	SourceId        = "sourceId"
	SourceTopic     = "sourceTopic"
	SourceNode      = "sourceNode"
	SourceFrom      = "sourceFrom"
	SourceQos       = "sourceQos"
	SourceTimestamp = "sourceTimestamp"
	SourcePayload   = "payload"
)

// MessageProvider 抽象出来的 消息提供者
type MessageProvider interface {

	// BatchSend 批量发送
	BatchSend(messages []*exhook_v2.Message)

	// SingleSend 单条发送
	SingleSend(message *exhook_v2.Message)
}
