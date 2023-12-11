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

type ProposalExecutor[T any] interface {
	Execute(props []*proposal.Proposal[T]) error
}

type MessageHandler[T any] interface {
	HandleMessage(m *message.Message[T]) (*proposal.Proposal[T], error)
}

// EVMChain is struct that aggregates all data required for
type EVMChain[T any] struct {
	listener       EventListener
	executor       ProposalExecutor[T]
	messageHandler MessageHandler[T]

	domainID   uint8
	startBlock *big.Int

	logger zerolog.Logger
}

func NewEVMChain[T any](listener EventListener, messageHandler MessageHandler[T], executor ProposalExecutor[T], domainID uint8, startBlock *big.Int) *EVMChain[T] {
	return &EVMChain[T]{
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
func (c *EVMChain[T]) PollEvents(ctx context.Context) {
	c.logger.Info().Str("startBlock", c.startBlock.String()).Msg("Polling Blocks...")
	go c.listener.ListenToEvents(ctx, c.startBlock)
}

func (c *EVMChain[T]) ReceiveMessage(m *message.Message[T]) (*proposal.Proposal[T], error) {
	return c.messageHandler.HandleMessage(m)
}

func (c *EVMChain[T]) Write(props []*proposal.Proposal[T]) error {
	err := c.executor.Execute(props)
	if err != nil {
		c.logger.Err(err).Msgf("error writing proposals %+v on network %d", props, c.DomainID())
		return err
	}

	return nil
}

func (c *EVMChain[T]) DomainID() uint8 {
	return c.domainID
}
