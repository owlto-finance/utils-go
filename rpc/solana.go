package rpc

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/owlto-finance/utils-go/loader"
	"github.com/owlto-finance/utils-go/log"
	sol "github.com/owlto-finance/utils-go/txn/solana"
	"github.com/owlto-finance/utils-go/util"
)

type SolanaRpc struct {
	tokenInfoMgr *loader.TokenInfoManager
	chainInfo    *loader.ChainInfo
}

func NewSolanaRpc(chainInfo *loader.ChainInfo) *SolanaRpc {
	return &SolanaRpc{
		chainInfo:    chainInfo,
		tokenInfoMgr: loader.NewTokenInfoManager(nil, nil),
	}
}

func (w *SolanaRpc) GetClient() *rpc.Client {
	return w.chainInfo.Client.(*rpc.Client)
}

func (w *SolanaRpc) GetAccount(ctx context.Context, ownerAddr string) (*rpc.Account, error) {
	ownerAddr = strings.TrimSpace(ownerAddr)

	ownerpk, err := solana.PublicKeyFromBase58(ownerAddr)
	if err != nil {
		return nil, err
	}

	rsp, err := w.GetClient().GetAccountInfo(
		ctx,
		ownerpk,
	)

	if err != nil {
		return nil, err
	} else {
		return rsp.Value, nil
	}
}

func (w *SolanaRpc) GetTokenInfo(ctx context.Context, tokenAddr string) (string, int32, error) {
	if util.IsHexStringZero(tokenAddr) {
		return "SOL", 9, nil
	}
	tokenInfo, ok := w.tokenInfoMgr.GetByChainNameTokenAddr(w.chainInfo.Name, tokenAddr)
	if ok {
		return tokenInfo.TokenName, tokenInfo.Decimals, nil
	}

	mintpk, err := solana.PublicKeyFromBase58(tokenAddr)
	if err != nil {
		return "", 0, err
	}

	metapk, _, err := solana.FindTokenMetadataAddress(mintpk)
	if err != nil {
		return "", 0, err
	}
	rsp, err := w.GetClient().GetAccountInfo(
		ctx,
		metapk,
	)
	if err != nil {
		return "", 0, err
	}

	meta, err := token_metadata.MetadataDeserialize(rsp.GetBinary())
	if err != nil {
		return "", 0, err
	}

	rsp, err = w.GetClient().GetAccountInfo(
		ctx,
		mintpk,
	)
	if err != nil {
		return "", 0, err
	}
	var mintAccount token.Mint
	decoder := bin.NewBorshDecoder(rsp.GetBinary())
	err = mintAccount.UnmarshalWithDecoder(decoder)
	if err != nil {
		return "", 0, err
	}

	w.tokenInfoMgr.AddToken(w.chainInfo.Name, meta.Data.Symbol, tokenAddr, int32(mintAccount.Decimals))
	return meta.Data.Symbol, int32(mintAccount.Decimals), nil

}

func (w *SolanaRpc) GetSplAccount(ctx context.Context, ownerAddr string, tokenAddr string) (*token.Account, error) {
	ownerAddr = strings.TrimSpace(ownerAddr)
	tokenAddr = strings.TrimSpace(tokenAddr)

	ownerpk, err := solana.PublicKeyFromBase58(ownerAddr)
	if err != nil {
		return nil, err
	}
	mintpk, err := solana.PublicKeyFromBase58(tokenAddr)
	if err != nil {
		return nil, err
	}

	ownerAta, err := sol.GetAtaFromPk(ownerpk, mintpk)
	if err != nil {
		return nil, err
	}

	rsp, err := w.GetClient().GetAccountInfo(
		ctx,
		ownerAta,
	)
	if err != nil {
		return nil, err
	}

	var tokenAccount token.Account
	decoder := bin.NewBorshDecoder(rsp.GetBinary())
	err = tokenAccount.UnmarshalWithDecoder(decoder)
	if err != nil {
		return nil, err
	} else {
		return &tokenAccount, nil
	}
}

func (w *SolanaRpc) GetBalance(ctx context.Context, ownerAddr string, tokenAddr string) (*big.Int, error) {
	ownerAddr = strings.TrimSpace(ownerAddr)
	tokenAddr = strings.TrimSpace(tokenAddr)

	if util.IsHexStringZero(tokenAddr) {
		account, err := w.GetAccount(ctx, ownerAddr)
		if err != nil {
			if err == rpc.ErrNotFound {
				return big.NewInt(0), nil
			}
			return nil, err
		}
		return big.NewInt(int64(account.Lamports)), nil
	} else {
		sqlAccount, err := w.GetSplAccount(ctx, ownerAddr, tokenAddr)
		if err != nil {
			if err == rpc.ErrNotFound {
				return big.NewInt(0), nil
			}
			return nil, err
		}
		return big.NewInt(int64(sqlAccount.Amount)), nil
	}
}

func (w *SolanaRpc) GetAllowance(ctx context.Context, ownerAddr string, tokenAddr string, spenderAddr string) (*big.Int, error) {
	sqlAccount, err := w.GetSplAccount(ctx, ownerAddr, tokenAddr)
	if err != nil {
		if err == rpc.ErrNotFound {
			return big.NewInt(0), nil
		}
		return nil, err
	} else {
		return big.NewInt(int64(sqlAccount.DelegatedAmount)), nil
	}
}

func (w *SolanaRpc) IsTxSuccess(ctx context.Context, hash string) (bool, int64, error) {
	sig, err := solana.SignatureFromBase58(hash)
	if err != nil {
		return false, 0, err
	}

	receipt, err := w.GetClient().GetTransaction(ctx, sig, nil)
	if err != nil {
		return false, 0, err
	}
	if receipt == nil {
		return false, 0, fmt.Errorf("get receipt failed")
	}
	return receipt.Meta.Err == nil, int64(receipt.Slot), nil
}

func (w *SolanaRpc) Client() interface{} {
	return w.chainInfo.Client
}

func (w *SolanaRpc) Backend() int32 {
	return 3
}

func (w *SolanaRpc) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	blockNumber, err := w.GetClient().GetSlot(
		context.TODO(),
		rpc.CommitmentFinalized,
	)

	if err != nil {
		log.Errorf("%v get latest block number error %v", w.chainInfo.Name, err)
		return 0, err
	}
	return int64(blockNumber), nil

}
