package impl

import (
	"context"
	"go_emqx_exhook/emqx.io/grpc/exhook"
	"log"
)

type HookProviderServerImpl struct {
	exhook.UnimplementedHookProviderServer
	SourceTopics []string
	Receive      func(msg *exhook.Message)
}

// OnProviderLoaded 注册钩子加载,开启钩子服务,onProviderLoaded中目前包含所有的钩子服务，可以将用的放开注释，只需要在本类中实现需要的方法即可
func (h *HookProviderServerImpl) OnProviderLoaded(ctx context.Context, request *exhook.ProviderLoadedRequest) (*exhook.LoadedResponse, error) {
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
	hooks := []*exhook.HookSpec{
		&exhook.HookSpec{
			Name:   "message.publish",
			Topics: h.SourceTopics,
		},
	}
	return &exhook.LoadedResponse{Hooks: hooks}, nil
}

// OnProviderUnloaded 关闭钩子服务
func (h *HookProviderServerImpl) OnProviderUnloaded(ctx context.Context, request *exhook.ProviderUnloadedRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

// OnClientConnect 客户端连接
func (h *HookProviderServerImpl) OnClientConnect(ctx context.Context, request *exhook.ClientConnectRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnClientConnack(ctx context.Context, request *exhook.ClientConnackRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

// OnClientConnected 客户端连接成功
func (h *HookProviderServerImpl) OnClientConnected(ctx context.Context, request *exhook.ClientConnectedRequest) (*exhook.EmptySuccess, error) {
	log.Printf("OnClientConnected: \n%v\n%v\n", request.GetClientinfo(), request.GetMeta())
	return &exhook.EmptySuccess{}, nil
}

// OnClientDisconnected 客户端断开连接
func (h *HookProviderServerImpl) OnClientDisconnected(ctx context.Context, request *exhook.ClientDisconnectedRequest) (*exhook.EmptySuccess, error) {
	log.Printf("OnClientDisconnected: \n%v\n%v\n%v\n", request.GetClientinfo(), request.GetMeta(), request.GetReason())
	return &exhook.EmptySuccess{}, nil
}

// OnClientAuthenticate 认证，单独开启认证功能后，该钩子服务失效
func (h *HookProviderServerImpl) OnClientAuthenticate(ctx context.Context, request *exhook.ClientAuthenticateRequest) (*exhook.ValuedResponse, error) {
	return &exhook.ValuedResponse{
		Type:  exhook.ValuedResponse_CONTINUE,
		Value: &exhook.ValuedResponse_BoolResult{BoolResult: true},
	}, nil
}

func (h *HookProviderServerImpl) OnClientAuthorize(ctx context.Context, request *exhook.ClientAuthorizeRequest) (*exhook.ValuedResponse, error) {
	return &exhook.ValuedResponse{
		Type:  exhook.ValuedResponse_CONTINUE,
		Value: &exhook.ValuedResponse_BoolResult{BoolResult: true},
	}, nil
}

func (h *HookProviderServerImpl) OnClientSubscribe(ctx context.Context, request *exhook.ClientSubscribeRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnClientUnsubscribe(ctx context.Context, request *exhook.ClientUnsubscribeRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnSessionCreated(ctx context.Context, request *exhook.SessionCreatedRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnSessionSubscribed(ctx context.Context, request *exhook.SessionSubscribedRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnSessionUnsubscribed(ctx context.Context, request *exhook.SessionUnsubscribedRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnSessionResumed(ctx context.Context, request *exhook.SessionResumedRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnSessionDiscarded(ctx context.Context, request *exhook.SessionDiscardedRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnSessionTakenover(ctx context.Context, request *exhook.SessionTakenoverRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnSessionTerminated(ctx context.Context, request *exhook.SessionTerminatedRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

// OnMessagePublish 收到消息处理
func (h *HookProviderServerImpl) OnMessagePublish(ctx context.Context, request *exhook.MessagePublishRequest) (*exhook.ValuedResponse, error) {
	h.Receive(request.GetMessage())
	return &exhook.ValuedResponse{
		Type:  exhook.ValuedResponse_CONTINUE,
		Value: &exhook.ValuedResponse_Message{Message: request.GetMessage()},
	}, nil
}

func (h *HookProviderServerImpl) OnMessageDelivered(ctx context.Context, request *exhook.MessageDeliveredRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnMessageDropped(ctx context.Context, request *exhook.MessageDroppedRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}

func (h *HookProviderServerImpl) OnMessageAcked(ctx context.Context, request *exhook.MessageAckedRequest) (*exhook.EmptySuccess, error) {
	return &exhook.EmptySuccess{}, nil
}
