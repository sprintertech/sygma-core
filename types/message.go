package types

type MessageType string
type Message struct {
	Source      uint8       // Source where message was initiated
	Destination uint8       // Destination chain of message
	Data        interface{} // Data associated with the message
	Type        MessageType // Message type
}

func NewMessage(
	source uint8,
	destination uint8,
	data interface{},
	msgType MessageType,
) *Message {
	return &Message{
		source,
		destination,
		data,
		msgType,
	}
}
