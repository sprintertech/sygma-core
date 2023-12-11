package proposal

type ProposalType string
type Proposal[T any] struct {
	Source      uint8
	Destination uint8
	Data        T
	Type        ProposalType
}

func NewProposal[T any](source, destination uint8, data T, propType ProposalType) *Proposal[T] {
	return &Proposal[T]{
		Source:      source,
		Destination: destination,
		Data:        data,
		Type:        propType,
	}
}
