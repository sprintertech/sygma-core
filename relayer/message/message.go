package message

type MessageType string
type Message struct {
	Source      uint8       // Source where message was initiated
	Destination uint8       // Destination chain of message
	Data        interface{} // Data associated with the message
	ID          string      // ID is used to track and identify message across networks
	Type        MessageType // Message type
}

func NewMessage(source, destination uint8, data interface{}, id string, msgType MessageType) *Message {
	return &Message{
		Source:      source,
		Destination: destination,
		Data:        data,
		Type:        msgType,
		ID:          id,
	}
}
