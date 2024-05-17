package loader

import (
	"database/sql"
	"strings"

	"github.com/owlto-finance/utils-go/alert"
)

type SrcTx struct {
	ChainId           int32
	TxHash            string
	Sender            string
	Receiver          string
	TargetAddress     sql.NullString
	Token             string
	Value             string
	DstChainid        sql.NullInt32
	IsTestnet         sql.NullInt32
	TxTimestamp       int32
	SrcTokenName      sql.NullString
	SrcTokenDecimal   int32
	IsCctp            int32
	SrcNonce          int32
	ThirdpartyChannel int32
}

type SrcTxManager struct {
	db      *sql.DB
	alerter alert.Alerter
	//mutex   *sync.RWMutex
}

func NewSrcTxManager(db *sql.DB, alerter alert.Alerter) *SrcTxManager {
	return &SrcTxManager{
		db:      db,
		alerter: alerter,
		//mutex:   &sync.RWMutex{},
	}
}

func (mgr *SrcTxManager) IsSrcTxExist(chainId int32, txHash string) bool {
	var id int64
	err := mgr.db.QueryRow("SELECT id FROM t_src_transaction where chainid = ? and tx_hash = ? ", chainId, strings.TrimSpace(txHash)).Scan(&id)
	return err == nil
}

func (mgr *SrcTxManager) SetResult(txHash string, isInvalid int32, isVerified int32) error {
	_, err := mgr.db.Exec("update t_src_transaction set is_invalid = ?, is_verified = ? where tx_hash = ? ", isInvalid, isVerified, txHash)
	if err != nil {
		mgr.alerter.AlertText("update t_transfer is_invalid error :", err)
		return err
	}
	return nil
}

func (mgr *SrcTxManager) Save(tx *SrcTx) error {
	tx.TxHash = strings.TrimSpace(tx.TxHash)
	tx.Sender = strings.TrimSpace(tx.Sender)
	tx.Receiver = strings.TrimSpace(tx.Receiver)
	tx.Token = strings.TrimSpace(tx.Token)
	tx.Value = strings.TrimSpace(tx.Value)
	tx.TargetAddress.String = strings.TrimSpace(tx.TargetAddress.String)
	tx.SrcTokenName.String = strings.TrimSpace(tx.SrcTokenName.String)

	query := `INSERT IGNORE INTO t_src_transaction (chainid, tx_hash, sender, receiver, target_address, token, value, dst_chainid, is_testnet, tx_timestamp, src_token_name, src_token_decimal, is_cctp, src_nonce, thirdparty_channel)
              VALUES (?, ?, ?, ?, ?, ?, ? , ?, ?, ?, ?, ?, ?, ?, ?)`

	// Execute the SQL statement with tx data
	_, err := mgr.db.Exec(query, tx.ChainId, tx.TxHash, tx.Sender, tx.Receiver, tx.TargetAddress, tx.Token, tx.Value, tx.DstChainid, tx.IsTestnet, tx.TxTimestamp, tx.SrcTokenName, tx.SrcTokenDecimal, tx.IsCctp, tx.SrcNonce, tx.ThirdpartyChannel)
	if err != nil {
		mgr.alerter.AlertText("failed to insert src transaction", err)
		return err
	}

	return nil

}
