package mqtt

type Props struct {
	PayloadFormatIndicator          uint8             // 载荷格式说明
	MessageExpiryInterval           uint32            // 消息过期时间
	ContentType                     string            // 内容类型
	ResponseTopic                   string            // 响应主题
	CorrelationData                 []byte            // 相关数据
	SubscriptionIdentifier          uint64            // 定义标识符
	SessionExpiryInterval           uint32            // 会话过期间隔
	AssignedClientIdentifier        string            // 分配客户标识符
	ServerKeepAlive                 uint16            // 服务端保活时间
	AuthenticationMethod            string            // 认证方法
	AuthenticationData              []byte            // 认证数据
	RequestProblemInformation       bool              // 请求问题信息
	WillDelayInterval               uint32            // 遗嘱延时间隔
	RequestResponseInformation      bool              // 请求响应信息
	ResponseInformation             string            // 请求信息
	ServerReference                 string            // 服务端参考
	ReasonString                    string            // 原因字符串
	ReceiveMaximum                  uint16            // 接收最大数量
	TopicAliasMaximum               uint16            // 主题别名最大长度
	TopicAlias                      uint16            // 主题别名
	MaximumQoS                      uint8             // 最大 QoS
	RetainAvailable                 bool              // 保留属性可用性
	UserProperty                    map[string]string // 用户属性
	MaximumPacketSize               uint32            // 最大报文长度
	WildcardSubscriptionAvailable   bool              // 通配符订阅可用性
	SubscriptionIdentifierAvailable bool              // 订阅标识符可用性
	SharedSubscriptionsAvailable    bool              // 共享订阅可用性
}


func (p *Props) ToConnect() *ConnectProps {
	return &ConnectProps{
		SessionExpiryInterval:      p.SessionExpiryInterval,
		AuthenticationMethod:       p.AuthenticationMethod,
		AuthenticationData:         p.AuthenticationData,
		RequestProblemInformation:  p.RequestProblemInformation,
		RequestResponseInformation: p.RequestResponseInformation,
		ReceiveMaximum:             p.ReceiveMaximum,
		TopicAliasMaximum:          p.TopicAliasMaximum,
		UserProperty:               p.UserProperty,
		MaximumPacketSize:          p.MaximumPacketSize,
	}
}

func (p *Props) ToWill() *WillProps {
	return &WillProps{
		PayloadFormatIndicator: p.PayloadFormatIndicator,
		MessageExpiryInterval:  p.MessageExpiryInterval,
		ContentType:            p.ContentType,
		ResponseTopic:          p.ResponseTopic,
		CorrelationData:        p.CorrelationData,
		WillDelayInterval:      p.WillDelayInterval,
		UserProperty:           p.UserProperty,
	}
}
