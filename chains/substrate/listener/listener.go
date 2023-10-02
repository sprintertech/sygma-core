// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"context"
	"math/big"
	"time"

	"github.com/ChainSafe/sygma-core/store"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventHandler interface {
	HandleEvents(evts []*parser.Event) error
}
type ChainConnection interface {
	UpdateMetatdata() error
	GetHeaderLatest() (*types.Header, error)
	GetBlockHash(blockNumber uint64) (types.Hash, error)
	GetBlockEvents(hash types.Hash) ([]*parser.Event, error)
	GetFinalizedHead() (types.Hash, error)
	GetBlock(blockHash types.Hash) (*types.SignedBlock, error)
}

type SubstrateListener struct {
	conn ChainConnection

	blockstore store.BlockStore

	eventHandlers []EventHandler

	blockRetryInterval time.Duration
	blockInterval      *big.Int
	domainID           uint8

	log zerolog.Logger
}

func NewSubstrateListener(connection ChainConnection, blockstore store.BlockStore, eventHandlers []EventHandler, domainID uint8, blockRetryInterval time.Duration, blockInterval *big.Int) *SubstrateListener {
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

				evts, err := l.fetchEvents(startBlock, endBlock)
				if err != nil {
					l.log.Err(err).Msgf("Failed fetching events for block range %s-%s", startBlock, endBlock)
					time.Sleep(l.blockRetryInterval)
					continue
				}

				for _, handler := range l.eventHandlers {
					err := handler.HandleEvents(evts)
					if err != nil {
						l.log.Error().Err(err).Msg("Error handling substrate events")
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

func (l *SubstrateListener) fetchEvents(startBlock *big.Int, endBlock *big.Int) ([]*parser.Event, error) {
	l.log.Debug().Msgf("Fetching substrate events for block range %s-%s", startBlock, endBlock)

	evts := make([]*parser.Event, 0)
	for i := new(big.Int).Set(startBlock); i.Cmp(endBlock) == -1; i.Add(i, big.NewInt(1)) {
		hash, err := l.conn.GetBlockHash(i.Uint64())
		if err != nil {
			return nil, err
		}

		evt, err := l.conn.GetBlockEvents(hash)
		if err != nil {
			return nil, err
		}
		evts = append(evts, evt...)

	}

	return evts, nil
}
