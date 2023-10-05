// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package evm

import (
	"context"
	"math/big"

	"github.com/ChainSafe/sygma-core/store"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventListener interface {
	ListenToEvents(ctx context.Context, startBlock *big.Int)
}

type ProposalExecutor interface {
	Execute(props []*types.Proposal) error
}

type MessageHandler interface {
	HandleMessage(m *types.Message) (*types.Proposal, error)
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener       EventListener
	executor       ProposalExecutor
	messageHandler MessageHandler

	blockstore *store.BlockStore

	domainID   uint8
	startBlock *big.Int

	logger zerolog.Logger
}

func NewEVMChain(listener EventListener, messageHandler MessageHandler, executor ProposalExecutor, blockstore *store.BlockStore, domainID uint8, startBlock *big.Int) *EVMChain {
	return &EVMChain{
		listener:       listener,
		executor:       executor,
		blockstore:     blockstore,
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

func (c *EVMChain) ReceiveMessage(m *types.Message) (*types.Proposal, error) {
	return c.messageHandler.HandleMessage(m)
}

func (c *EVMChain) Write(props []*types.Proposal) error {
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
