package rpc

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/owlto-finance/utils-go/loader"
	"github.com/owlto-finance/utils-go/network"
	"github.com/owlto-finance/utils-go/util"
)

type ZksliteRpc struct {
	chainInfo *loader.ChainInfo
}

func NewZksliteRpc(chainInfo *loader.ChainInfo) *ZksliteRpc {
	return &ZksliteRpc{
		chainInfo: chainInfo,
	}
}

func (w *ZksliteRpc) Client() interface{} {
	return w.chainInfo.Client
}

func (w *ZksliteRpc) Backend() int32 {
	return 1
}

func (w *ZksliteRpc) GetTokenInfo(ctx context.Context, tokenAddr string) (string, int32, error) {
	return "", 0, fmt.Errorf("no impl")
}

func (w *ZksliteRpc) GetAllowance(ctx context.Context, ownerAddr string, tokenAddr string, spenderAddr string) (*big.Int, error) {
	return big.NewInt(0), fmt.Errorf("not impl")
}

func (w *ZksliteRpc) IsLastCharSlash(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[len(s)-1] == '/'
}

func (w *ZksliteRpc) GetBalance(ctx context.Context, ownerAddr string, tokenAddr string) (*big.Int, error) {
	ownerAddr = strings.TrimSpace(ownerAddr)
	tokenAddr = strings.TrimSpace(tokenAddr)

	var data map[string]interface{}
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "account_info",
		"params":  []string{ownerAddr},
	}

	var url = w.chainInfo.RpcEndPoint
	if w.IsLastCharSlash(url) {
		url += "jsrpc"
	} else {
		url += "/jsrpc"
	}

	if util.IsHexStringZero(tokenAddr) || tokenAddr == "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" || tokenAddr == "0xdAC17F958D2ee523a2206206994597C13D831ec7" {
		err := network.Request(url, request, &data)
		if err != nil {
			return big.NewInt(0), err
		}

		if result, ok := data["result"].(map[string]interface{}); ok {
			if committed, ok := result["committed"].(map[string]interface{}); ok {
				if balances, ok := committed["balances"].(map[string]interface{}); ok {
					if util.IsHexStringZero(tokenAddr) {
						if eth, ok := balances["ETH"].(string); ok {
							value, ok := big.NewInt(0).SetString(eth, 10)
							if ok {
								return value, nil
							}
						}
					} else if tokenAddr == "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" {
						if usdc, ok := balances["USDC"].(string); ok {
							value, ok := big.NewInt(0).SetString(usdc, 10)
							if ok {
								return value, nil
							}
						}
					} else if tokenAddr == "0xdAC17F958D2ee523a2206206994597C13D831ec7" {
						if usdt, ok := balances["USDT"].(string); ok {
							value, ok := big.NewInt(0).SetString(usdt, 10)
							if ok {
								return value, nil
							}
						}
					}
				}
			}
		}
		return big.NewInt(0), nil
	} else {
		return big.NewInt(0), fmt.Errorf("not impl")
	}
}

func (w *ZksliteRpc) IsTxSuccess(ctx context.Context, hash string) (bool, int64, error) {
	return false, 0, fmt.Errorf("not impl")
}

func (w *ZksliteRpc) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	return 0, fmt.Errorf("not impl")
}
