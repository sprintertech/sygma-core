package calls

import (
	"context"
	"math/big"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxFabric func(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrices []*big.Int, data []byte) (client.CommonTransaction, error)

type ContractChecker interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
}

type ContractCaller interface {
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
}

type GasPricer interface {
	// make priority a pointer to uint8 to pass nil into all GasPrice functions (instead of magic numbers)
	GasPrice(priority *uint8) ([]*big.Int, error)
}

type ClientDispatcher interface {
	TxReceipt(h common.Hash) (*types.Receipt, error)
	SignAndSendTransaction(ctx context.Context, tx client.CommonTransaction) (common.Hash, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	GetTransactionByHash(h common.Hash) (tx *types.Transaction, isPending bool, err error)
	UnsafeNonce() (*big.Int, error)
	LockNonce()
	UnlockNonce()
	UnsafeIncreaseNonce() error
	From() common.Address
}

type ContractCallerDispatcher interface {
	ContractCaller
	ClientDispatcher
	ContractChecker
}
