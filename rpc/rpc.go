package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/owlto-finance/utils-go/loader"
)

type Rpc interface {
	Client() interface{}
	Backend() int32
	GetLatestBlockNumber(ctx context.Context) (int64, error)
	IsTxSuccess(ctx context.Context, hash string) (bool, int64, error)
	GetAllowance(ctx context.Context, ownerAddr string, tokenAddr string, spenderAddr string) (*big.Int, error)
	GetBalance(ctx context.Context, ownerAddr string, tokenAddr string) (*big.Int, error)
}

func GetRpc(chainInfo *loader.ChainInfo) (Rpc, error) {
	if chainInfo.Backend == 1 {
		return NewEvmRpc(chainInfo), nil
	} else if chainInfo.Backend == 2 {
		return NewStarknetRpc(chainInfo), nil
	} else if chainInfo.Backend == 3 {
		return NewSolanaRpc(chainInfo), nil
	} else if chainInfo.Backend == 4 {
		return NewBitcoinRpc(chainInfo), nil
	} else if chainInfo.Backend == 5 {
		return NewZksliteRpc(chainInfo), nil
	}
	return nil, fmt.Errorf("unsupport backend %v", chainInfo.Backend)
}
