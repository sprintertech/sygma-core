package proposal

type ProposalType string
type Proposal struct {
	Source      uint8
	Destination uint8
	Data        interface{}
	Type        ProposalType
}

func NewProposal(source uint8, destination uint8, data interface{}, propType ProposalType) *Proposal {
	return &Proposal{
		Source:      source,
		Destination: destination,
		Data:        data,
		Type:        propType,
	}
}
