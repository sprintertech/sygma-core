package message

import "github.com/ChainSafe/sygma-core/relayer/proposal"

type Handler interface {
	HandleMessage(m *Message) (*proposal.Proposal, error)
}

type MessageHandler struct {
	handlers map[MessageType]Handler
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

// HandlerMessage calls associated handler for that message type and returns a proposal to be submitted on-chain
func (h *MessageHandler) HandleMessage(m *Message) (*proposal.Proposal, error) {
	return h.handlers[m.Type].HandleMessage(m)
}

// RegisterMessageHandler registers an message handler by associating a handler to a message type
func (mh *MessageHandler) RegisterMessageHandler(t MessageType, h Handler) {
	mh.handlers[t] = h
}
