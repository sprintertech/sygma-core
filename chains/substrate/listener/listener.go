// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"context"
	"math/big"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventHandler interface {
	HandleEvents(startBlock *big.Int, endBlock *big.Int) error
}

type ChainConnection interface {
	UpdateMetatdata() error
	GetHeaderLatest() (*types.Header, error)
	GetBlockHash(blockNumber uint64) (types.Hash, error)
	GetBlockEvents(hash types.Hash) ([]*parser.Event, error)
	GetFinalizedHead() (types.Hash, error)
	GetBlock(blockHash types.Hash) (*types.SignedBlock, error)
}

type BlockStorer interface {
	StoreBlock(block *big.Int, domainID uint8) error
}

type SubstrateListener struct {
	conn          ChainConnection
	blockstore    BlockStorer
	eventHandlers []EventHandler

	blockRetryInterval time.Duration
	blockInterval      *big.Int
	domainID           uint8

	log zerolog.Logger
}

func NewSubstrateListener(connection ChainConnection, blockstore BlockStorer, eventHandlers []EventHandler, domainID uint8, blockRetryInterval time.Duration, blockInterval *big.Int) *SubstrateListener {
	return &SubstrateListener{
		log:                log.With().Uint8("domainID", domainID).Logger(),
		domainID:           domainID,
		conn:               connection,
		blockstore:         blockstore,
		eventHandlers:      eventHandlers,
		blockRetryInterval: blockRetryInterval,
		blockInterval:      blockInterval,
	}
}

func (l *SubstrateListener) ListenToEvents(ctx context.Context, startBlock *big.Int) {
	endBlock := big.NewInt(0)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				hash, err := l.conn.GetFinalizedHead()
				if err != nil {
					l.log.Error().Err(err).Msg("Failed to fetch finalized header")
					time.Sleep(l.blockRetryInterval)
					continue
				}
				head, err := l.conn.GetBlock(hash)
				if err != nil {
					l.log.Error().Err(err).Msg("Failed to fetch block")
					time.Sleep(l.blockRetryInterval)
					continue
				}

				if startBlock == nil {
					startBlock = big.NewInt(int64(head.Block.Header.Number))
				}
				endBlock.Add(startBlock, l.blockInterval)

				// Sleep if finalized is less then current block
				if big.NewInt(int64(head.Block.Header.Number)).Cmp(endBlock) == -1 {
					time.Sleep(l.blockRetryInterval)
					continue
				}

				for _, handler := range l.eventHandlers {
					err := handler.HandleEvents(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)))
					if err != nil {
						l.log.Warn().Err(err).Msg("Error handling substrate events")
						continue
					}
				}

				err = l.blockstore.StoreBlock(startBlock, l.domainID)
				if err != nil {
					l.log.Error().Str("block", startBlock.String()).Err(err).Msg("Failed to write latest block to blockstore")
				}
				startBlock.Add(startBlock, l.blockInterval)
			}
		}
	}()
}
