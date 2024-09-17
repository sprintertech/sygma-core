package message

type MessageType string
type Message struct {
	Source      uint8       // Source where message was initiated
	Destination uint8       // Destination chain of message
	Data        interface{} // Data associated with the message
	ID          string      // ID is used to track and identify message across networks
	Type        MessageType // Message type
	ErrChn      chan error  // ErrChn is used to share errors that happen on the destination handler
}

func NewMessage(source, destination uint8, data interface{}, id string, msgType MessageType, errChn chan error) *Message {
	return &Message{
		Source:      source,
		Destination: destination,
		Data:        data,
		Type:        msgType,
		ID:          id,
		ErrChn:      errChn,
	}
}
