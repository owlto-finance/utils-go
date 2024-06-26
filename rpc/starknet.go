package rpc

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/owlto-finance/utils-go/loader"
	"github.com/owlto-finance/utils-go/log"
)

type StarknetRpc struct {
	chainInfo *loader.ChainInfo
}

func NewStarknetRpc(chainInfo *loader.ChainInfo) *StarknetRpc {
	return &StarknetRpc{
		chainInfo: chainInfo,
	}
}

func (w *StarknetRpc) GetClient() *rpc.Provider {
	return w.chainInfo.Client.(*rpc.Provider)
}

func (w *StarknetRpc) Client() interface{} {
	return w.chainInfo.Client
}

func (w *StarknetRpc) GetTokenInfo(ctx context.Context, tokenAddr string) (string, int32, error) {
	return "", 0, fmt.Errorf("no impl")
}

func (w *StarknetRpc) GetBalance(ctx context.Context, ownerAddr string, tokenAddr string) (*big.Int, error) {
	ownerAddr = strings.TrimSpace(ownerAddr)
	tokenAddr = strings.TrimSpace(tokenAddr)

	token, err := utils.HexToFelt(tokenAddr)
	if err != nil {
		return nil, err
	}
	owner, err := utils.HexToFelt(ownerAddr)
	if err != nil {
		return nil, err
	}
	tx := rpc.FunctionCall{
		ContractAddress:    token,
		EntryPointSelector: utils.GetSelectorFromNameFelt("balanceOf"),
		Calldata:           []*felt.Felt{owner},
	}
	rsp, err := w.GetClient().Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return nil, err
	}
	if len(rsp) > 0 {
		return rsp[0].BigInt(new(big.Int)), nil
	} else {
		return big.NewInt(0), nil
	}
}

func (w *StarknetRpc) GetAllowance(ctx context.Context, ownerAddr string, tokenAddr string, spenderAddr string) (*big.Int, error) {
	return nil, fmt.Errorf("starknet get allowance unsupport")
}

func (w *StarknetRpc) Backend() int32 {
	return 2
}

func (w *StarknetRpc) IsTxSuccess(ctx context.Context, hash string) (bool, int64, error) {
	bhash, err := hexutil.Decode(hash)
	if err != nil {
		return false, 0, err
	}
	status, err := w.GetClient().GetTransactionStatus(ctx, new(felt.Felt).SetBytes(bhash))
	if err != nil {
		return false, 0, err
	}
	if status == nil {
		return false, 0, fmt.Errorf("get status failed")
	}

	switch status.FinalityStatus {
	case rpc.TxnStatus_Rejected:
		return false, 0, nil
	case rpc.TxnStatus_Received:
		// ignore this status, wait for L2 confirmation
		return false, 0, fmt.Errorf("not complete: %v", status.FinalityStatus)
	case rpc.TxnStatus_Accepted_On_L2, rpc.TxnStatus_Accepted_On_L1:
		return true, 0, nil //status.ExecutionStatus == rpc.TxnExecutionStatusSUCCEEDED, nil
	default:
		return false, 0, fmt.Errorf("unknown tx status: %v", status.FinalityStatus)
	}
}

func (w *StarknetRpc) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	blockNumber, err := w.GetClient().BlockNumber(ctx)
	if err != nil {
		log.Errorf("%v get latest block number error %v", w.chainInfo.Name, err)
		return 0, err
	}
	return int64(blockNumber), nil
}
