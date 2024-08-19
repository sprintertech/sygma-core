// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package connection

import (
	"math/big"
	"sync"

	"github.com/sygmaprotocol/go-substrate-rpc-client/v4/client"
	"github.com/sygmaprotocol/go-substrate-rpc-client/v4/registry/parser"
	"github.com/sygmaprotocol/go-substrate-rpc-client/v4/registry/retriever"
	"github.com/sygmaprotocol/go-substrate-rpc-client/v4/registry/state"

	"github.com/sygmaprotocol/go-substrate-rpc-client/v4/rpc"
	"github.com/sygmaprotocol/go-substrate-rpc-client/v4/rpc/chain"
	"github.com/sygmaprotocol/go-substrate-rpc-client/v4/types"
)

type CheckMetadataModeEnabledConnection struct {
	chain.Chain
	client.Client
	*rpc.RPC
	meta        types.Metadata // Latest chain metadata
	metaLock    sync.RWMutex   // Lock metadata for updates, allows concurrent reads
	GenesisHash types.Hash     // Chain genesis hash
}

func NewCheckMetadataModeEnabledConnection(url string) (*CheckMetadataModeEnabledConnection, error) {
	client, err := client.Connect(url)
	if err != nil {
		return nil, err
	}
	rpc, err := rpc.NewRPC(client)
	if err != nil {
		return nil, err
	}

	meta, err := rpc.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	genesisHash, err := rpc.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}

	return &CheckMetadataModeEnabledConnection{
		meta: *meta,

		RPC:         rpc,
		Chain:       rpc.Chain,
		Client:      client,
		GenesisHash: types.Hash(genesisHash),
	}, nil
}

func (c *CheckMetadataModeEnabledConnection) GetMetadata() (meta types.Metadata) {
	c.metaLock.RLock()
	meta = c.meta
	c.metaLock.RUnlock()
	return meta
}

func (c *CheckMetadataModeEnabledConnection) UpdateMetatdata() error {
	c.metaLock.Lock()
	meta, err := c.RPC.State.GetMetadataLatest()
	if err != nil {
		c.metaLock.Unlock()
		return err
	}
	c.meta = *meta
	c.metaLock.Unlock()
	return nil
}

func (c *CheckMetadataModeEnabledConnection) GetBlockEvents(hash types.Hash) ([]*parser.Event, error) {
	provider := state.NewEventProvider(c.State)
	eventRetriever, err := retriever.NewDefaultEventRetriever(provider, c.State)
	if err != nil {
		return nil, err
	}

	evts, err := eventRetriever.GetEvents(hash)
	if err != nil {
		return nil, err
	}
	return evts, nil
}

func (c *CheckMetadataModeEnabledConnection) FetchEvents(startBlock, endBlock *big.Int) ([]*parser.Event, error) {
	evts := make([]*parser.Event, 0)
	for i := new(big.Int).Set(startBlock); i.Cmp(endBlock) <= 0; i.Add(i, big.NewInt(1)) {
		hash, err := c.GetBlockHash(i.Uint64())
		if err != nil {
			return nil, err
		}

		evt, err := c.GetBlockEvents(hash)
		if err != nil {
			return nil, err
		}
		evts = append(evts, evt...)
	}
	return evts, nil
}
