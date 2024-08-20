package substrate

import (
	"context"
	"fmt"
	"math/big"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sygmaprotocol/sygma-core/relayer/message"
	"github.com/sygmaprotocol/sygma-core/relayer/proposal"
)

type ProposalExecutor interface {
	Execute(props []*proposal.Proposal) error
}

type MessageHandler interface {
	HandleMessage(m *message.Message) (*proposal.Proposal, error)
}

type EventListener interface {
	ListenToEvents(ctx context.Context, startBlock *big.Int)
}

type SubstrateChain struct {
	listener       EventListener
	messageHandler MessageHandler
	executor       ProposalExecutor

	domainID   uint8
	startBlock *big.Int

	logger zerolog.Logger
}

func NewSubstrateChain(listener EventListener, messageHandler MessageHandler, executor ProposalExecutor, domainID uint8, startBlock *big.Int) *SubstrateChain {
	return &SubstrateChain{
		listener:       listener,
		messageHandler: messageHandler,
		executor:       executor,
		domainID:       domainID,
		startBlock:     startBlock,
		logger:         log.With().Uint8("domainID", domainID).Logger()}
}

// PollEvents is the goroutine that polls blocks and searches Deposit events in them.
// Events are then sent to eventsChan.
func (c *SubstrateChain) PollEvents(ctx context.Context) {
	if c.listener == nil {
		return
	}

	c.logger.Info().Str("startBlock", c.startBlock.String()).Msg("Polling Blocks...")
	go c.listener.ListenToEvents(ctx, c.startBlock)
}

func (c *SubstrateChain) ReceiveMessage(m *message.Message) (*proposal.Proposal, error) {
	if c.messageHandler == nil {
		return nil, fmt.Errorf("message handler not configured")
	}

	return c.messageHandler.HandleMessage(m)
}

func (c *SubstrateChain) Write(props []*proposal.Proposal) error {
	if c.executor == nil {
		return fmt.Errorf("executor not configured")
	}

	err := c.executor.Execute(props)
	if err != nil {
		c.logger.Err(err).Str("messageID", props[0].MessageID).Msgf("error writing proposals %+v on network %d", props, c.DomainID())
		return err
	}

	return nil
}

func (c *SubstrateChain) DomainID() uint8 {
	return c.domainID
}
