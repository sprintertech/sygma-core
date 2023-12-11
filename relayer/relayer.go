// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package relayer

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sygmaprotocol/sygma-core/relayer/message"
	"github.com/sygmaprotocol/sygma-core/relayer/proposal"
)

type RelayedChain[T any] interface {
	// PollEvents starts listening for on-chain events
	PollEvents(ctx context.Context)
	// ReceiveMessage accepts the message from the source chain and converts it into
	// a Proposal to be submitted on-chain
	ReceiveMessage(m *message.Message[T]) (*proposal.Proposal[T], error)
	// Write submits proposals on-chain.
	// If multiple proposals submitted they are expected to be able to be batched.
	Write(proposals []*proposal.Proposal[T]) error
	DomainID() uint8
}

func NewRelayer[T any](chains map[uint8]RelayedChain[T]) *Relayer[T] {
	return &Relayer[T]{relayedChains: chains}
}

type Relayer[T any] struct {
	relayedChains map[uint8]RelayedChain[T]
}

// Start function starts polling events for each chain and listens to cross-chain messages.
// If an array of messages is sent to the channel they are expected to be to the same destination and
// able to be handled in batches.
func (r *Relayer[T]) Start(ctx context.Context, msgChan chan []*message.Message[T]) {
	log.Info().Msgf("Starting relayer")

	for _, c := range r.relayedChains {
		log.Debug().Msgf("Starting chain %v", c.DomainID())
		go c.PollEvents(ctx)
	}

	for {
		select {
		case m := <-msgChan:
			go r.route(m)
			continue
		case <-ctx.Done():
			return
		}
	}
}

// Route function routes the messages to the destination chain.
func (r *Relayer[T]) route(msgs []*message.Message[T]) {
	destChain, ok := r.relayedChains[msgs[0].Destination]
	if !ok {
		log.Error().Uint8("domainID", destChain.DomainID()).Msgf("No chain registered for destination domain")
		return
	}

	props := make([]*proposal.Proposal[T], 0)
	for _, m := range msgs {
		prop, err := destChain.ReceiveMessage(m)
		if err != nil {
			log.Err(err).Uint8("domainID", destChain.DomainID()).Msgf("Failed receiving message %+v", m)
			continue
		}
		if prop != nil {
			props = append(props, prop)
		}
	}
	if len(props) == 0 {
		return
	}

	err := destChain.Write(props)
	if err != nil {
		log.Err(err).Uint8("domainID", destChain.DomainID()).Msgf("Failed writing message")
		return
	}
}
