// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package chains

import (
	"github.com/ChainSafe/sygma-core/relayer/message"
	"github.com/ChainSafe/sygma-core/types"
)

func NewProposal(source, destination uint8, depositNonce uint64, resourceId types.ResourceID, data []byte, metadata message.Metadata) *Proposal {
	return &Proposal{
		OriginDomainID: source,
		DepositNonce:   depositNonce,
		ResourceID:     resourceId,
		Destination:    destination,
		Data:           data,
		Metadata:       metadata,
	}
}

type Proposal struct {
	OriginDomainID uint8
	DepositNonce   uint64
	ResourceID     types.ResourceID
	Data           []byte
	Destination    uint8
	Metadata       message.Metadata
}
