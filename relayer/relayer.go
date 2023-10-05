// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package relayer

import (
	"context"

	"github.com/ChainSafe/sygma-core/relayer/message"
	"github.com/ChainSafe/sygma-core/relayer/proposal"
	"github.com/rs/zerolog/log"
)

type RelayedChain interface {
	PollEvents(ctx context.Context)
	ReceiveMessage(m *message.Message) (*proposal.Proposal, error)
	Write(proposals []*proposal.Proposal) error
	DomainID() uint8
}

func NewRelayer(chains map[uint8]RelayedChain) *Relayer {
	return &Relayer{relayedChains: chains}
}

type Relayer struct {
	relayedChains map[uint8]RelayedChain
}

// Start function starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (r *Relayer) Start(ctx context.Context, msgChan chan []*message.Message) {
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

// Route function runs destination writer by mapping DestinationID from message to registered writer.
func (r *Relayer) route(msgs []*message.Message) {
	destChain, ok := r.relayedChains[msgs[0].Destination]
	if !ok {
		log.Error().Uint8("domainID", destChain.DomainID()).Msgf("no resolver for destID %v to send message registered", msgs[0].Destination)
		return
	}

	props := make([]*proposal.Proposal, 0)
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
