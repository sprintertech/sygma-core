package relayer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ChainSafe/sygma-core/mock"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/stretchr/testify/suite"
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
	s.mockRelayedChain.EXPECT().ReceiveMessages(gomock.Any()).Return(make([]*types.Proposal, 0), fmt.Errorf("error"))
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	msgChan := make(chan []*types.Message, 1)
	msgChan <- []*types.Message{
		{Destination: 1},
	}
	relayer.Start(ctx, msgChan)
	time.Sleep(time.Millisecond * 100)
}

func (s *RouteTestSuite) TestReceiveMessageFails() {
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1)).Times(1)
	s.mockRelayedChain.EXPECT().ReceiveMessages(gomock.Any()).Return(make([]*types.Proposal, 0), fmt.Errorf("error"))
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	relayer.route([]*types.Message{
		{Destination: 1},
	})
}

func (s *RouteTestSuite) TestAvoidWriteWithoutProposals() {
	s.mockRelayedChain.EXPECT().ReceiveMessages(gomock.Any()).Return(make([]*types.Proposal, 0), nil)
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	relayer.route([]*types.Message{
		{Destination: 1},
	})
}

func (s *RouteTestSuite) TestWriteFails() {
	props := make([]*types.Proposal, 1)
	prop := &types.Proposal{}
	props[0] = prop
	s.mockRelayedChain.EXPECT().ReceiveMessages(gomock.Any()).Return(props, nil)
	s.mockRelayedChain.EXPECT().Write(props).Return(fmt.Errorf("error"))
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1)).Times(1)
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	relayer.route([]*types.Message{
		{Destination: 1},
	})
}

func (s *RouteTestSuite) TestWritesToChain() {
	props := make([]*types.Proposal, 1)
	prop := &types.Proposal{}
	props[0] = prop
	s.mockRelayedChain.EXPECT().ReceiveMessages(gomock.Any()).Return(props, nil)
	s.mockRelayedChain.EXPECT().Write(props).Return(nil)
	chains := make(map[uint8]RelayedChain)
	chains[1] = s.mockRelayedChain
	relayer := NewRelayer(
		chains,
	)

	relayer.route([]*types.Message{
		{Destination: 1},
	})
}
