package relayer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/sygmaprotocol/sygma-core/mock"
	"github.com/sygmaprotocol/sygma-core/relayer/message"
	"github.com/sygmaprotocol/sygma-core/relayer/proposal"
	"go.uber.org/mock/gomock"
)

type RouteTestSuite struct {
	suite.Suite
	mockRelayedChain *mock.MockRelayedChain
}

func TestRunRouteTestSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

func (s *RouteTestSuite) SetupSuite()    {}
func (s *RouteTestSuite) TearDownSuite() {}
func (s *RouteTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockRelayedChain = mock.NewMockRelayedChain(gomockController)
}
func (s *RouteTestSuite) TearDownTest() {}

func (s *RouteTestSuite) TestStartListensOnChannel() {
	ctx, cancel := context.WithCancel(context.TODO())

	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1))
	s.mockRelayedChain.EXPECT().PollEvents(gomock.Any())
	s.mockRelayedChain.EXPECT().DomainID().DoAndReturn(func() uint8 {
		cancel()
		return 1
	})
	s.mockRelayedChain.EXPECT().ReceiveMessage(gomock.Any()).Return(nil, fmt.Errorf("error"))
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	msgChan := make(chan []*message.Message, 1)
	msgChan <- []*message.Message{
		{Destination: 1},
	}
	relayer.Start(ctx, msgChan)
	time.Sleep(time.Millisecond * 100)
}

func (s *RouteTestSuite) TestReceiveMessageFails() {
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1)).Times(1)
	s.mockRelayedChain.EXPECT().ReceiveMessage(gomock.Any()).Return(nil, fmt.Errorf("error"))
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	relayer.route([]*message.Message{
		{Destination: 1},
	})
}

func (s *RouteTestSuite) TestAvoidWriteWithoutProposals() {
	s.mockRelayedChain.EXPECT().ReceiveMessage(gomock.Any()).Return(nil, nil)
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	relayer.route([]*message.Message{
		{Destination: 1},
	})
}

func (s *RouteTestSuite) TestWriteFails() {
	props := make([]*proposal.Proposal, 1)
	prop := &proposal.Proposal{}
	props[0] = prop
	s.mockRelayedChain.EXPECT().ReceiveMessage(gomock.Any()).Return(prop, nil)
	s.mockRelayedChain.EXPECT().Write(props).Return(fmt.Errorf("error"))
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1)).Times(1)
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	relayer.route([]*message.Message{
		{Destination: 1},
	})
}

func (s *RouteTestSuite) TestWritesToChain() {
	props := make([]*proposal.Proposal, 1)
	prop := &proposal.Proposal{}
	props[0] = prop
	s.mockRelayedChain.EXPECT().ReceiveMessage(gomock.Any()).Return(prop, nil)
	s.mockRelayedChain.EXPECT().Write(props).Return(nil)
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	relayer.route([]*message.Message{
		{Destination: 1},
	})
}
