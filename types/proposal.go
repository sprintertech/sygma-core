package types

import (
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/types"
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
