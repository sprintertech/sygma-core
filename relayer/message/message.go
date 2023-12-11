package message

type MessageType string
type Message[T any] struct {
	Source      uint8       // Source where message was initiated
	Destination uint8       // Destination chain of message
	Data        T           // Data associated with the message
	Type        MessageType // Message type
}

func NewMessage[T any](
	source uint8,
	destination uint8,
	data T,
	msgType MessageType,
) *Message[T] {
	return &Message[T]{
		source,
		destination,
		data,
		msgType,
	}
}
