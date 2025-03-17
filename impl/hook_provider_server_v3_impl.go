package impl

import (
	"context"
	"go_emqx_exhook/emqx.io/grpc/exhook_v3"
	"log"
)

type HookProviderServerV3Impl struct {
	exhook_v3.UnimplementedHookProviderServer
	SourceTopics []string
	Callback     func(request *exhook_v3.MessagePublishRequest)
}

// OnProviderLoaded 注册钩子加载,开启钩子服务,onProviderLoaded中目前包含所有的钩子服务，可以将用的放开注释，只需要在本类中实现需要的方法即可
func (h *HookProviderServerV3Impl) OnProviderLoaded(ctx context.Context, request *exhook_v3.ProviderLoadedRequest) (*exhook_v3.LoadedResponse, error) {
	log.Printf("OnProviderLoaded: %v \n", request.Broker)
	/*		名称					说明				执行时机
	client.connect			处理连接报文		服务端收到客户端的连接报文时
	client.connack			下发连接应答		服务端准备下发连接应答报文时
	client.connected		成功接入			客户端认证完成并成功接入系统后
	client.disconnected		连接断开			客户端连接层在准备关闭时
	client.authenticate		连接认证			执行完 client.connect 后
	client.authorize		发布订阅鉴权		执行 发布/订阅 操作前
	client.subscribe		订阅主题			收到订阅报文后，执行 client.authorize 鉴权前
	client.unsubscribe		取消订阅			收到取消订阅报文后
	session.created			会话创建			client.connected 执行完成，且创建新的会话后
	session.subscribed		会话订阅主题		完成订阅操作后
	session.unsubscribed	会话取消订阅		完成取消订阅操作后
	session.resumed			会话恢复			client.connected 执行完成，且成功恢复旧的会话信息后
	session.discarded		会话被移除		会话由于被移除而终止后
	session.takenover		会话被接管		会话由于被接管而终止后
	session.terminated		会话终止			会话由于其他原因被终止后
	message.publish			消息发布			服务端在发布（路由）消息前
	message.delivered		消息投递			消息准备投递到客户端前
	message.acked			消息回执			服务端在收到客户端发回的消息 ACK 后
	message.dropped			消息丢弃			发布出的消息被丢弃后
	*/
	log.Printf("subscribed topics : %v \n", h.SourceTopics)
	hooks := []*exhook_v3.HookSpec{
		&exhook_v3.HookSpec{
			Name:   "message.publish",
			Topics: h.SourceTopics,
		},
	}
	return &exhook_v3.LoadedResponse{Hooks: hooks}, nil
}

// OnProviderUnloaded 关闭钩子服务
func (h *HookProviderServerV3Impl) OnProviderUnloaded(ctx context.Context, request *exhook_v3.ProviderUnloadedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

// OnClientConnect 客户端连接
func (h *HookProviderServerV3Impl) OnClientConnect(ctx context.Context, request *exhook_v3.ClientConnectRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnClientConnack(ctx context.Context, request *exhook_v3.ClientConnackRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

// OnClientConnected 客户端连接成功
func (h *HookProviderServerV3Impl) OnClientConnected(ctx context.Context, request *exhook_v3.ClientConnectedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

// OnClientDisconnected 客户端断开连接
func (h *HookProviderServerV3Impl) OnClientDisconnected(ctx context.Context, request *exhook_v3.ClientDisconnectedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

// OnClientAuthenticate 认证，单独开启认证功能后，该钩子服务失效
func (h *HookProviderServerV3Impl) OnClientAuthenticate(ctx context.Context, request *exhook_v3.ClientAuthenticateRequest) (*exhook_v3.ValuedResponse, error) {
	return &exhook_v3.ValuedResponse{
		Type:  exhook_v3.ValuedResponse_CONTINUE,
		Value: &exhook_v3.ValuedResponse_BoolResult{BoolResult: true},
	}, nil
}

func (h *HookProviderServerV3Impl) OnClientAuthorize(ctx context.Context, request *exhook_v3.ClientAuthorizeRequest) (*exhook_v3.ValuedResponse, error) {
	return &exhook_v3.ValuedResponse{
		Type:  exhook_v3.ValuedResponse_CONTINUE,
		Value: &exhook_v3.ValuedResponse_BoolResult{BoolResult: true},
	}, nil
}

func (h *HookProviderServerV3Impl) OnClientSubscribe(ctx context.Context, request *exhook_v3.ClientSubscribeRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnClientUnsubscribe(ctx context.Context, request *exhook_v3.ClientUnsubscribeRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnSessionCreated(ctx context.Context, request *exhook_v3.SessionCreatedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnSessionSubscribed(ctx context.Context, request *exhook_v3.SessionSubscribedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnSessionUnsubscribed(ctx context.Context, request *exhook_v3.SessionUnsubscribedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnSessionResumed(ctx context.Context, request *exhook_v3.SessionResumedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnSessionDiscarded(ctx context.Context, request *exhook_v3.SessionDiscardedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnSessionTakenover(ctx context.Context, request *exhook_v3.SessionTakenoverRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnSessionTerminated(ctx context.Context, request *exhook_v3.SessionTerminatedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

// OnMessagePublish 收到消息处理
func (h *HookProviderServerV3Impl) OnMessagePublish(ctx context.Context, request *exhook_v3.MessagePublishRequest) (*exhook_v3.ValuedResponse, error) {
	// 消息发送到管道
	h.Callback(request)
	return &exhook_v3.ValuedResponse{
		Type:  exhook_v3.ValuedResponse_CONTINUE,
		Value: &exhook_v3.ValuedResponse_Message{Message: request.GetMessage()},
	}, nil
}

func (h *HookProviderServerV3Impl) OnMessageDelivered(ctx context.Context, request *exhook_v3.MessageDeliveredRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnMessageDropped(ctx context.Context, request *exhook_v3.MessageDroppedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}

func (h *HookProviderServerV3Impl) OnMessageAcked(ctx context.Context, request *exhook_v3.MessageAckedRequest) (*exhook_v3.EmptySuccess, error) {
	return &exhook_v3.EmptySuccess{}, nil
}
