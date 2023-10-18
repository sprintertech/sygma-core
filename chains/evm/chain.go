// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package evm

import (
	"context"
	"math/big"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sygmaprotocol/sygma-core/store"
	"github.com/sygmaprotocol/sygma-core/types"
)

type EventListener interface {
	ListenToEvents(ctx context.Context, startBlock *big.Int)
}

type ProposalExecutor interface {
	Execute(messages []*types.Message) error
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener   EventListener
	executor   ProposalExecutor
	blockstore *store.BlockStore

	domainID   uint8
	startBlock *big.Int

	logger zerolog.Logger
}

func NewEVMChain(listener EventListener, executor ProposalExecutor, blockstore *store.BlockStore, domainID uint8, startBlock *big.Int) *EVMChain {
	return &EVMChain{
		listener:   listener,
		executor:   executor,
		blockstore: blockstore,
		domainID:   domainID,
		startBlock: startBlock,
		logger:     log.With().Uint8("domainID", domainID).Logger(),
	}
}

// PollEvents is the goroutine that polls blocks and searches Deposit events in them.
// Events are then sent to eventsChan.
func (c *EVMChain) PollEvents(ctx context.Context) {
	c.logger.Info().Str("startBlock", c.startBlock.String()).Msg("Polling Blocks...")
	go c.listener.ListenToEvents(ctx, c.startBlock)
}

func (c *EVMChain) Write(msgs []*types.Message) error {
	err := c.executor.Execute(msgs)
	if err != nil {
		c.logger.Err(err).Msgf("error writing messages %+v on network %d", msgs, c.DomainID())
		return err
	}

	return nil
}

func (c *EVMChain) DomainID() uint8 {
	return c.domainID
}
