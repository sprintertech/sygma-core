package message

import (
	"fmt"

	"github.com/sygmaprotocol/sygma-core/relayer/proposal"
)

type Handler[T any] interface {
	HandleMessage(m *Message[T]) (*proposal.Proposal[T], error)
}

type MessageHandler[T any] struct {
	handlers map[MessageType]Handler[T]
}

func NewMessageHandler[T any]() *MessageHandler[T] {
	return &MessageHandler[T]{
		handlers: make(map[MessageType]Handler[T]),
	}
}

// HandlerMessage calls associated handler for that message type and returns a proposal to be submitted on-chain
func (h *MessageHandler[T]) HandleMessage(m *Message[T]) (*proposal.Proposal[T], error) {
	mh, ok := h.handlers[m.Type]
	if !ok {
		return nil, fmt.Errorf("no handler found for type %s", m.Type)
	}
	return mh.HandleMessage(m)
}

// RegisterMessageHandler registers a message handler by associating a handler to a message type
func (mh *MessageHandler[T]) RegisterMessageHandler(t MessageType, h Handler[T]) {
	mh.handlers[t] = h
}
