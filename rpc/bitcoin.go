package rpc

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	_ "github.com/gagliardetto/solana-go"
	"github.com/owlto-finance/utils-go/loader"
	"github.com/owlto-finance/utils-go/network"
	"github.com/owlto-finance/utils-go/util"
)

type BitcoinRpc struct {
	chainInfo *loader.ChainInfo
}

func NewBitcoinRpc(chainInfo *loader.ChainInfo) *BitcoinRpc {
	return &BitcoinRpc{
		chainInfo: chainInfo,
	}
}

func (w *BitcoinRpc) GetBalance(ctx context.Context, ownerAddr string, tokenAddr string) (*big.Int, error) {
	ownerAddr = strings.TrimSpace(ownerAddr)
	tokenAddr = strings.TrimSpace(tokenAddr)

	if util.IsHexStringZero(tokenAddr) {
		var data map[string]interface{}
		request := map[string]interface{}{
			"method": "bb_getaddress",
			"params": []any{ownerAddr, map[string]interface{}{
				"page":       1,
				"size":       1,
				"fromHeight": 0,
				"details":    "basic",
			}},
		}
		err := network.Request(w.chainInfo.RpcEndPoint, request, &data)
		if err != nil {
			return big.NewInt(0), err
		}

		if result, ok := data["result"].(map[string]interface{}); ok {
			if balance, ok := result["balance"].(string); ok {
				value, ok := big.NewInt(0).SetString(balance, 10)
				if ok {
					return value, nil
				}
			}
		}
		return big.NewInt(0), nil
	} else {
		return big.NewInt(0), fmt.Errorf("not impl")
	}
}

func (w *BitcoinRpc) GetAllowance(ctx context.Context, ownerAddr string, tokenAddr string, spenderAddr string) (*big.Int, error) {
	return big.NewInt(0), fmt.Errorf("not impl")
}

func (w *BitcoinRpc) IsTxSuccess(ctx context.Context, hash string) (bool, int64, error) {
	return false, 0, fmt.Errorf("not impl")
}

func (w *BitcoinRpc) Client() interface{} {
	return w.chainInfo.Client
}

func (w *BitcoinRpc) Backend() int32 {
	return 4
}

func (w *BitcoinRpc) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	return 0, fmt.Errorf("not impl")
}
