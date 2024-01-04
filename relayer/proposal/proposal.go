package proposal

type ProposalType string
type Proposal struct {
	Source      uint8
	Destination uint8
	Data        interface{}
	Type        ProposalType
}

func NewProposal(source, destination uint8, data []byte, propType ProposalType) *Proposal {
	return &Proposal{
		Source:      source,
		Destination: destination,
		Data:        data,
		Type:        propType,
	}
}
