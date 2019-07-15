package mqtt

type (
	PropType uint8
	QoSLevel uint8
	ReasonCode uint8
	PacketType uint8
)

const (
	PayloadFormatIndicator         PropType = 0x01 // 载荷格式说明
	MessageExpiryInterval          PropType = 0x02 // 消息过期时间
	ContentType                    PropType = 0x03 // 内容类型
	ResponseTopic                  PropType = 0x08 // 响应主题
	CorrelationData                PropType = 0x09 // 相关数据
	SubscriptionIdentifier         PropType = 0x0B // 定义标识符
	SessionExpiryInterval          PropType = 0x11 // 会话过期间隔
	AssignedClientIdentifier       PropType = 0x12 // 分配客户标识符
	ServerKeepAlive                PropType = 0x13 // 服务端保活时间
	AuthenticationMethod           PropType = 0x15 // 认证方法
	AuthenticationData             PropType = 0x16 // 认证数据
	RequestProblemInformation      PropType = 0x17 // 请求问题信息
	WillDelayInterval              PropType = 0x18 // 遗嘱延时间隔
	RequestResponseInformation     PropType = 0x19 // 请求响应信息
	ResponseInformation            PropType = 0x1A // 请求信息
	ServerReference                PropType = 0x1C // 服务端参考
	ReasonString                   PropType = 0x1F // 原因字符串
	ReceiveMaximum                 PropType = 0x21 // 接收最大数量
	TopicAliasMaximum              PropType = 0x22 // 接收最大数量
	TopicAlias                     PropType = 0x23 // 接收最大数量
	MaximumQoS                     PropType = 0x24 // 最大 QoS
	RetainAvalilable               PropType = 0x25 // 保留属性可用性
	UserProperty                   PropType = 0x26 // 用户属性
	MaximumPacketSize              PropType = 0x27 // 最大报文长度
	WildcardSubscriptionAvailable  PropType = 0x28 // 通配符订阅可用性
	SubscrptionIdentifierAvailable PropType = 0x29 // 订阅标识符可用性
	SharedSubscriptionsAvaliable   PropType = 0x2A // 共享订阅可用性
)

const (
	Success                             ReasonCode = 0x00
	NormalDisconnection                 ReasonCode = 0x00
	GrantedQoS0                         ReasonCode = 0x00
	GrantedQoS1                         ReasonCode = 0x01
	GrantedQoS2                         ReasonCode = 0x02
	DisconnectWithWillMessage           ReasonCode = 0x04
	NoMatchingSubscribers               ReasonCode = 0x10
	NoSubscriptionExisted               ReasonCode = 0x11
	ContinueAuthentication              ReasonCode = 0x18
	ReAuthenticate                      ReasonCode = 0x19
	UnspecifiedError                    ReasonCode = 0x80
	MalformedPacket                     ReasonCode = 0x81
	ProtocolError                       ReasonCode = 0x82
	ImplementationSpecificError         ReasonCode = 0x83
	UnsupportedProtocolVersion          ReasonCode = 0x84
	ClientIdentifierNotValid            ReasonCode = 0x85
	BadUsernameOrPassword               ReasonCode = 0x86
	NotAuthorized                       ReasonCode = 0x87
	ServerUnavailable                   ReasonCode = 0x88
	ServerBusy                          ReasonCode = 0x89
	Banned                              ReasonCode = 0x8A
	ServerShuttingDown                  ReasonCode = 0x8B
	BadAuthenticationMethod             ReasonCode = 0x8C
	KeepAliveTimeout                    ReasonCode = 0x8D
	SessionTakenOver                    ReasonCode = 0x8E
	TopicFilterInvalid                  ReasonCode = 0x8F
	TopicNameInvalid                    ReasonCode = 0x90
	PacketIdentifierInUse               ReasonCode = 0x91
	PacketIdentifierNotFound            ReasonCode = 0x92
	ReceiveMaximumExceeded              ReasonCode = 0x93
	TopicAliasInvalid                   ReasonCode = 0x94
	PacketTooLarge                      ReasonCode = 0x95
	MessageRateTooHigh                  ReasonCode = 0x96
	QuotaExceeded                       ReasonCode = 0x97
	AdministrativeAction                ReasonCode = 0x98
	PayloadFormatInvalid                ReasonCode = 0x99
	RetianlNotSupported                 ReasonCode = 0x9A
	QoSNotSupported                     ReasonCode = 0x9B
	UseAnotherServer                    ReasonCode = 0x9C
	ServerMoved                         ReasonCode = 0x9D
	SharedSubscriptionsNotSupported     ReasonCode = 0x9E
	ConnectionRateExceeded              ReasonCode = 0x9F
	MaximumConnectionTime               ReasonCode = 0xA0
	SubscriptionIdentifiersNotSupported ReasonCode = 0xA1
	WildcardSubscriptionsNotSupported   ReasonCode = 0xA2
)

const (
	QoS0 QoSLevel = iota
	QoS1
	QoS2
)

const (
	_           PacketType = iota
	CONNECT      // 客户端请求连接服务端
	CONNACK      // 连接报文确认
	PUBLISH      // 发布消息
	PUBACK       // QoS 1 消息发布收到确认
	PUBREC       // 发布收到（保证交付第)
	PUBREL       // 发布释放（保证交付第二步）
	PUBCOMP      // QoS 2 消息发布完成（保证交付第三步）
	SUBSCRIBE    // 客户端订阅请求
	SUBACK       // 订阅请求报文确认
	UNSUBSCRIBE  // 客户端取消订阅请求
	UNSUBACK     // 取消订阅报文确认
	PINGREQ      // 心跳请求
	PINGRESP     // 心跳响应
	DISCONNECT   // 断开连接通知
	AUTH         // 认证信息交换
)

func (p PropType) Byte() byte {
	return byte(p)
}