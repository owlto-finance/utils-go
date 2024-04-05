package sol

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"
)

type SolanaAccount struct {
	PublicKey  solana.PublicKey `json:"public_key"`
	IsWritable bool             `json:"is_writable"`
	IsSigner   bool             `json:"is_signer"`
}

type SolanaInstruction struct {
	ProgramId solana.PublicKey `json:"program_id"`
	Accounts  []SolanaAccount  `json:"accounts"`
	Data      hexutil.Bytes    `json:"data"`
}

type SolanaBody struct {
	Instructions []SolanaInstruction `json:"instructions"`
}

func GetAta(addr string, mint string) (solana.PublicKey, error) {
	pk, err := solana.PublicKeyFromBase58(addr)
	if err != nil {
		return solana.PublicKey{}, err
	}

	mintpk, err := solana.PublicKeyFromBase58(mint)
	if err != nil {
		return solana.PublicKey{}, err
	}

	ata, _, err := solana.FindAssociatedTokenAddress(pk, mintpk)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return ata, nil

}

func GetAtaFromPk(pk solana.PublicKey, mintpk solana.PublicKey) (solana.PublicKey, error) {

	ata, _, err := solana.FindAssociatedTokenAddress(pk, mintpk)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return ata, nil

}

func ToBody(insts []solana.Instruction) ([]byte, error) {
	body := SolanaBody{
		Instructions: make([]SolanaInstruction, 0, len(insts)),
	}

	for _, inst := range insts {
		data, err := inst.Data()
		if err != nil {
			return nil, err
		}

		maccs := make([]SolanaAccount, 0, len(inst.Accounts()))
		for _, acc := range inst.Accounts() {
			macc := SolanaAccount{
				PublicKey:  acc.PublicKey,
				IsWritable: acc.IsWritable,
				IsSigner:   acc.IsSigner,
			}
			maccs = append(maccs, macc)
		}

		minst := SolanaInstruction{
			ProgramId: inst.ProgramID(),
			Accounts:  maccs,
			Data:      data,
		}

		body.Instructions = append(body.Instructions, minst)
	}

	// Marshal the map to a JSON string
	return json.Marshal(body)
}
