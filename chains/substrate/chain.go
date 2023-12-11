package substrate

import (
	"context"
	"math/big"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sygmaprotocol/sygma-core/relayer/message"
	"github.com/sygmaprotocol/sygma-core/relayer/proposal"
)

type ProposalExecutor[T any] interface {
	Execute(props []*proposal.Proposal[T]) error
}

type MessageHandler[T any] interface {
	HandleMessage(m *message.Message[T]) (*proposal.Proposal[T], error)
}

type EventListener interface {
	ListenToEvents(ctx context.Context, startBlock *big.Int)
}

type SubstrateChain[T any] struct {
	listener       EventListener
	messageHandler MessageHandler[T]
	executor       ProposalExecutor[T]

	domainID   uint8
	startBlock *big.Int

	logger zerolog.Logger
}

func NewSubstrateChain[T any](listener EventListener, messageHandler MessageHandler[T], executor ProposalExecutor[T], domainID uint8, startBlock *big.Int) *SubstrateChain[T] {
	return &SubstrateChain[T]{
		listener:   listener,
		executor:   executor,
		domainID:   domainID,
		startBlock: startBlock,
		logger:     log.With().Uint8("domainID", domainID).Logger()}
}

// PollEvents is the goroutine that polls blocks and searches Deposit events in them.
// Events are then sent to eventsChan.
func (c *SubstrateChain[T]) PollEvents(ctx context.Context) {
	c.logger.Info().Str("startBlock", c.startBlock.String()).Msg("Polling Blocks...")
	go c.listener.ListenToEvents(ctx, c.startBlock)
}

func (c *SubstrateChain[T]) ReceiveMessage(m *message.Message[T]) (*proposal.Proposal[T], error) {
	return c.messageHandler.HandleMessage(m)
}

func (c *SubstrateChain[T]) Write(props []*proposal.Proposal[T]) error {
	err := c.executor.Execute(props)
	if err != nil {
		c.logger.Err(err).Msgf("error writing proposals %+v on network %d", props, c.DomainID())
		return err
	}

	return nil
}

func (c *SubstrateChain[T]) DomainID() uint8 {
	return c.domainID
}
