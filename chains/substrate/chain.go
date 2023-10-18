package substrate

import (
	"context"
	"math/big"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sygmaprotocol/sygma-core/chains/substrate/client"
	"github.com/sygmaprotocol/sygma-core/store"
	"github.com/sygmaprotocol/sygma-core/types"
)

type BatchProposalExecutor interface {
	Execute(msgs []*types.Message) error
}

type SubstrateChain struct {
	client *client.SubstrateClient

	listener EventListener
	executor BatchProposalExecutor

	blockstore *store.BlockStore

	domainID   uint8
	startBlock *big.Int

	logger zerolog.Logger
}

type EventListener interface {
	ListenToEvents(ctx context.Context, startBlock *big.Int)
}

func NewSubstrateChain(client *client.SubstrateClient, listener EventListener, blockstore *store.BlockStore, executor BatchProposalExecutor, domainID uint8, startBlock *big.Int) *SubstrateChain {
	return &SubstrateChain{
		client:     client,
		listener:   listener,
		blockstore: blockstore,
		executor:   executor,
		logger:     log.With().Uint8("domainID", domainID).Logger()}
}

// PollEvents is the goroutine that polls blocks and searches Deposit events in them.
// Events are then sent to eventsChan.
func (c *SubstrateChain) PollEvents(ctx context.Context) {
	c.logger.Info().Str("startBlock", c.startBlock.String()).Msg("Polling Blocks...")
	go c.listener.ListenToEvents(ctx, c.startBlock)
}

func (c *SubstrateChain) Write(msgs []*types.Message) error {
	err := c.executor.Execute(msgs)
	if err != nil {
		c.logger.Err(err).Msgf("error writing messages %+v on network %d", msgs, c.DomainID())
		return err
	}

	return nil
}

func (c *SubstrateChain) DomainID() uint8 {
	return c.domainID
}
