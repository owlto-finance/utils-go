package sol

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

type SolanaKeypair struct {
	PublicKey  solana.PublicKey  `json:"public_key"`
	PrivateKey solana.PrivateKey `json:"private_key"`
}

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
	Instructions []SolanaInstruction                        `json:"instructions"`
	Keypairs     []SolanaKeypair                            `json:"keypairs"`
	LookupTables map[solana.PublicKey]solana.PublicKeySlice `json:"lookup_tables"`
}

func (body *SolanaBody) AddInstructions(insts []solana.Instruction) error {
	for _, inst := range insts {
		data, err := inst.Data()
		if err != nil {
			return err
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
	return nil
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

func TransferBody(senderAddr string, receiverAddr string, amount *big.Int) ([]byte, error) {
	senderAddr = strings.TrimSpace(senderAddr)
	receiverAddr = strings.TrimSpace(receiverAddr)

	senderpk, err := solana.PublicKeyFromBase58(senderAddr)
	if err != nil {
		return nil, err
	}
	receiverpk, err := solana.PublicKeyFromBase58(receiverAddr)
	if err != nil {
		return nil, err
	}

	inst := system.NewTransferInstruction(
		amount.Uint64(),
		senderpk,
		receiverpk,
	).Build()

	return ToBody([]solana.Instruction{inst}, nil)

}

func ToBody(insts []solana.Instruction, keypairs []SolanaKeypair) ([]byte, error) {
	body, err := ToSolanaBody(insts, keypairs, nil)
	if err != nil {
		return nil, err
	}
	return json.Marshal(body)
}

func ToSolanaBody(insts []solana.Instruction, keypairs []SolanaKeypair, lookupTables map[solana.PublicKey]solana.PublicKeySlice) (*SolanaBody, error) {
	body := SolanaBody{
		Instructions: make([]SolanaInstruction, 0, len(insts)),
	}
	if keypairs != nil {
		body.Keypairs = keypairs
	}

	if lookupTables != nil {
		body.LookupTables = lookupTables
	}

	body.AddInstructions(insts)
	return &body, nil

}
