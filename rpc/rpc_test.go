package rpc

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/owlto-finance/utils-go/loader"
)

func TestSol(t *testing.T) {
	t.Log("test sol...")
	solRpc := NewSolanaRpc(&loader.ChainInfo{Name: "SolanaMainnet", Client: rpc.New("https://api.mainnet-beta.solana.com")})
	t.Log(solRpc.GetTokenInfo(context.TODO(), "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "J8qZijXxrypJin5Y27qcTvNjmd5ybF44NJdDKCSkXxWv"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "JC2iABaZUucksCEHZ95NCxb2BaBVhu8g4efaTNtyVsYL"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "J5tzd1ww1V1qrgDUQHVCGqpmpbnEnjzGs9LAqJxwkNde"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "zxTtD4MMnEAgHMvXmfgPCyMY61ivxX5zwu12hTSqLoA"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "J8qZijXxrypJin5Y27qcTvNjmd5ybF44NJdDKCSkXxWv"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "JC2iABaZUucksCEHZ95NCxb2BaBVhu8g4efaTNtyVsYL"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "J5tzd1ww1V1qrgDUQHVCGqpmpbnEnjzGs9LAqJxwkNde"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "zxTtD4MMnEAgHMvXmfgPCyMY61ivxX5zwu12hTSqLoA"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"))
	t.Log(solRpc.GetTokenInfo(context.TODO(), "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"))

}
