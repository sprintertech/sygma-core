// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package executor

import (
	"fmt"

	"github.com/sygmaprotocol/sygma-core/types"

	"github.com/rs/zerolog/log"
)

type Handlers map[types.TransferType]MessageHandlerFunc
type MessageHandlerFunc func(m *types.Message) (*types.Proposal, error)

type SubstrateMessageHandler struct {
	handlers Handlers
}

// NewSubstrateMessageHandler creates an instance of SubstrateMessageHandler that contains
// message handler functions for converting deposit message into a chain specific
// proposal
func NewSubstrateMessageHandler() *SubstrateMessageHandler {
	return &SubstrateMessageHandler{
		handlers: make(map[types.TransferType]MessageHandlerFunc),
	}
}

func (mh *SubstrateMessageHandler) HandleMessage(m *types.Message) (*types.Proposal, error) {
	// Based on handler that was registered on BridgeContract
	handleMessage, err := mh.matchTransferTypeHandlerFunc(m.Type)
	if err != nil {
		return nil, err
	}
	log.Info().Str("type", string(m.Type)).Uint8("src", m.Source).Uint8("dst", m.Destination).Uint64("nonce", m.DepositNonce).Str("resourceID", fmt.Sprintf("%x", m.ResourceId)).Msg("Handling new message")
	prop, err := handleMessage(m)
	if err != nil {
		return nil, err
	}
	return prop, nil
}

func (mh *SubstrateMessageHandler) matchTransferTypeHandlerFunc(transferType types.TransferType) (MessageHandlerFunc, error) {
	h, ok := mh.handlers[transferType]
	if !ok {
		return nil, fmt.Errorf("no corresponding message handler for this transfer type %s exists", transferType)
	}
	return h, nil
}

// RegisterEventHandler registers an message handler by associating a handler function to a specified transfer type
func (mh *SubstrateMessageHandler) RegisterMessageHandler(transferType types.TransferType, handler MessageHandlerFunc) {
	if transferType == "" {
		return
	}

	log.Info().Msgf("Registered message handler for transfer type %s", transferType)

	mh.handlers[transferType] = handler
}
