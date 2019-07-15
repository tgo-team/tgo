package mqtt


type PublishProps struct {
	PayloadFormatIndicator uint8
	MessageExpiryInterval  uint32
	ContentType            string
	ResponseTopic          string
	CorrelationData        []byte
	SubscriptionIdentifier uint64
	TopicAlias             uint16
	UserProperty           map[string]string
}

func (pp *PublishProps) ToProps() *Props {
	return &Props{
		PayloadFormatIndicator:pp.PayloadFormatIndicator,
		MessageExpiryInterval:pp.MessageExpiryInterval,
		ContentType:pp.ContentType,
		ResponseTopic:pp.ResponseTopic,
		CorrelationData:pp.CorrelationData,
		SubscriptionIdentifier:pp.SubscriptionIdentifier,
		TopicAlias:pp.TopicAlias,
		UserProperty:pp.UserProperty,

	}
}

type PublishVariableHeader struct {
	TopicName string
	Identifier  uint16

	props *PublishProps

}

func NewPublishVariableHeader(identifier uint16,topicName string,props *PublishProps) *PublishVariableHeader {

	return &PublishVariableHeader{props:props,TopicName:topicName,Identifier:identifier}
}


func (p *PublishVariableHeader) GetProps() *Props {
	return p.props.ToProps()
}