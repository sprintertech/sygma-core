// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package evm

import (
	"context"
	"math/big"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sygmaprotocol/sygma-core/relayer/message"
	"github.com/sygmaprotocol/sygma-core/relayer/proposal"
)

type EventListener interface {
	ListenToEvents(ctx context.Context, startBlock *big.Int)
}

type ProposalExecutor interface {
	Execute(props []*proposal.Proposal[any]) error
}

type MessageHandler interface {
	HandleMessage(m *message.Message[any]) (*proposal.Proposal[any], error)
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener       EventListener
	executor       ProposalExecutor
	messageHandler MessageHandler

	domainID   uint8
	startBlock *big.Int

	logger zerolog.Logger
}

func NewEVMChain(listener EventListener, messageHandler MessageHandler, executor ProposalExecutor, domainID uint8, startBlock *big.Int) *EVMChain {
	return &EVMChain{
		listener:       listener,
		executor:       executor,
		domainID:       domainID,
		startBlock:     startBlock,
		messageHandler: messageHandler,
		logger:         log.With().Uint8("domainID", domainID).Logger(),
	}
}

// PollEvents is the goroutine that polls blocks and searches Deposit events in them.
// Events are then sent to eventsChan.
func (c *EVMChain) PollEvents(ctx context.Context) {
	c.logger.Info().Str("startBlock", c.startBlock.String()).Msg("Polling Blocks...")
	go c.listener.ListenToEvents(ctx, c.startBlock)
}

func (c *EVMChain) ReceiveMessage(m *message.Message[any]) (*proposal.Proposal[any], error) {
	return c.messageHandler.HandleMessage(m)
}

func (c *EVMChain) Write(props []*proposal.Proposal[any]) error {
	err := c.executor.Execute(props)
	if err != nil {
		c.logger.Err(err).Msgf("error writing proposals %+v on network %d", props, c.DomainID())
		return err
	}

	return nil
}

func (c *EVMChain) DomainID() uint8 {
	return c.domainID
}
