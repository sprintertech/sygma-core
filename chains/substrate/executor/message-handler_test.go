// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package executor_test

import (
	"bytes"
	"errors"
	"math/big"
	"testing"
	"unsafe"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	substrateTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sygmaprotocol/sygma-core/chains/substrate/executor"
	"github.com/sygmaprotocol/sygma-core/types"

	"github.com/stretchr/testify/suite"
)

var SubstratePK = signature.KeyringPair{
	URI:       "//Alice",
	PublicKey: []byte{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d},
	Address:   "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
}

type MessageHandlerTestSuite struct {
	suite.Suite
}

func TestRunFungibleTransferHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(MessageHandlerTestSuite))
}

func (s *MessageHandlerTestSuite) TestSuccesfullyRegisterFungibleTransferMessageHandler() {
	recipientAddr := *(*[]substrateTypes.U8)(unsafe.Pointer(&SubstratePK.PublicKey))
	recipient := ConstructRecipientData(recipientAddr)

	messageData := &types.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         "fungible",
		Payload: []interface{}{
			[]byte{2}, // amount
			recipient,
		},
		Metadata: types.Metadata{},
	}

	invalidMessageData := &types.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         "nonFungible",
		Payload: []interface{}{
			[]byte{2}, // amount
			recipient,
		},
		Metadata: types.Metadata{},
	}

	depositMessageHandler := executor.NewSubstrateMessageHandler()
	// Register FungibleTransferMessageHandler function
	depositMessageHandler.RegisterMessageHandler("fungible", FungibleMessageHandler)
	prop1, err1 := depositMessageHandler.HandleMessage(messageData)
	s.Nil(err1)
	s.NotNil(prop1)

	// Use unregistered transfer type
	prop2, err2 := depositMessageHandler.HandleMessage(invalidMessageData)
	s.Nil(prop2)
	s.NotNil(err2)
}

func FungibleMessageHandler(m *types.Message) (*types.Proposal, error) {
	if len(m.Payload) != 2 {
		return nil, errors.New("malformed payload. Len  of payload should be 2")
	}
	amount, ok := m.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("wrong payload amount format")
	}
	recipient, ok := m.Payload[1].([]byte)
	if !ok {
		return nil, errors.New("wrong payload recipient format")
	}
	var data []byte
	data = append(data, common.LeftPadBytes(amount, 32)...) // amount (uint256)

	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	data = append(data, common.LeftPadBytes(recipientLen, 32)...)
	data = append(data, recipient...)
	return types.NewProposal(m.Source, m.Destination, m.DepositNonce, m.ResourceId, data, m.Metadata), nil
}

func ConstructRecipientData(recipient []substrateTypes.U8) []byte {
	rec := substrateTypes.MultiLocationV1{
		Parents: 0,
		Interior: substrateTypes.JunctionsV1{
			IsX1: true,
			X1: substrateTypes.JunctionV1{
				IsAccountID32: true,
				AccountID32NetworkID: substrateTypes.NetworkID{
					IsAny: true,
				},
				AccountID: recipient,
			},
		},
	}

	encodedRecipient := bytes.NewBuffer([]byte{})
	encoder := scale.NewEncoder(encodedRecipient)
	_ = rec.Encode(*encoder)

	recipientBytes := encodedRecipient.Bytes()
	var finalRecipient []byte

	// remove accountID size data
	// this is a fix because the substrate decoder is not able to parse the data with extra data
	// that represents size of the recipient byte array
	finalRecipient = append(finalRecipient, recipientBytes[:4]...)
	finalRecipient = append(finalRecipient, recipientBytes[5:]...)

	return finalRecipient
}
