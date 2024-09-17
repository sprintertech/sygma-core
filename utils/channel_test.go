package utils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/sygmaprotocol/sygma-core/utils"
)

type ChannelTestSuite struct {
	suite.Suite
}

func TestRunChannelTestSuite(t *testing.T) {
	suite.Run(t, new(ChannelTestSuite))
}

func (s *ChannelTestSuite) Test_TrySendError_NonBlocking() {
	utils.TrySendError(nil, fmt.Errorf("error"))

	errChn := make(chan error)
	utils.TrySendError(errChn, fmt.Errorf("error"))

	bufErrChn := make(chan error, 1)
	utils.TrySendError(bufErrChn, fmt.Errorf("error"))

	err := <-bufErrChn
	s.NotNil(err)
}
