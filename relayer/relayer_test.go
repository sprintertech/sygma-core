package relayer

import (
	"fmt"
	"testing"

	mock_relayer "github.com/ChainSafe/sygma-core/relayer/mock"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RouteTestSuite struct {
	suite.Suite
	mockRelayedChain *mock_relayer.MockRelayedChain
	mockMetrics      *mock_relayer.MockDepositMeter
}

func TestRunRouteTestSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

func (s *RouteTestSuite) SetupSuite()    {}
func (s *RouteTestSuite) TearDownSuite() {}
func (s *RouteTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockRelayedChain = mock_relayer.NewMockRelayedChain(gomockController)
	s.mockMetrics = mock_relayer.NewMockDepositMeter(gomockController)
}
func (s *RouteTestSuite) TearDownTest() {}

func (s *RouteTestSuite) TestLogsErrorIfDestinationDoesNotExist() {
	relayer := Relayer{
		metrics: s.mockMetrics,
	}

	relayer.route([]*types.Message{
		{},
	})
}

func (s *RouteTestSuite) TestLogsErrorIfMessageProcessorReturnsError() {
	s.mockMetrics.EXPECT().TrackDepositMessage(gomock.Any())
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1))
	relayer := NewRelayer(
		[]RelayedChain{},
		s.mockMetrics,
	)
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route([]*types.Message{
		{Destination: 1},
	})
}

func (s *RouteTestSuite) TestWriteFail() {
	s.mockMetrics.EXPECT().TrackDepositMessage(gomock.Any())
	s.mockMetrics.EXPECT().TrackExecutionError(gomock.Any())
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1)).Times(3)
	s.mockRelayedChain.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("error"))
	relayer := NewRelayer(
		[]RelayedChain{},
		s.mockMetrics,
	)
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route([]*types.Message{
		{Destination: 1},
	})
}

func (s *RouteTestSuite) TestWritesToDestChainIfMessageValid() {
	s.mockMetrics.EXPECT().TrackDepositMessage(gomock.Any())
	s.mockMetrics.EXPECT().TrackSuccessfulExecutionLatency(gomock.Any())
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1)).Times(2)
	s.mockRelayedChain.EXPECT().Write(gomock.Any())
	relayer := NewRelayer(
		[]RelayedChain{},
		s.mockMetrics,
	)
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route([]*types.Message{
		{Destination: 1},
	})
}
