package proposal

type ProposalType string
type Proposal struct {
	Source      uint8
	Destination uint8
	Data        interface{}
	Type        ProposalType
}
