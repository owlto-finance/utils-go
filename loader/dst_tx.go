package loader

import (
	"database/sql"

	"github.com/owlto-finance/utils-go/alert"
)

type DstTx struct {
	SrcAction         string
	SrcId             int64
	SrcVersion        int32
	Sender            int64
	Body              string
	FeeCap            sql.NullString
	TransferToken     sql.NullString
	TransferRecipient sql.NullString
	TransferAmount    sql.NullString
}

type TxGen struct {
	Id               int64
	Hash             string
	ConfirmedSuccess int8
}

type DstTxManager struct {
	db      *sql.DB
	alerter alert.Alerter
	//mutex   *sync.RWMutex
}

func NewDstTxManager(db *sql.DB, alerter alert.Alerter) *AccountManager {
	return &AccountManager{
		db:      db,
		alerter: alerter,
		//mutex:   &sync.RWMutex{},
	}
}

func (mgr *DstTxManager) GetDoneTxGenBySrc(srcId int64, action string, version int32) *TxGen {
	genId := mgr.GetDstTxConfirmGen(srcId, action, version)
	if genId == 0 {
		return nil
	}
	return mgr.GetDoneTxGen(genId)
}

func (mgr *DstTxManager) GetDoneTxGen(genId int64) *TxGen {
	var gen TxGen
	err := mgr.db.QueryRow("SELECT id,hash, confirmed_success FROM t_dst_transaction_gen where id = ? and confirmed_success is not null", genId).Scan(&gen.Id, &gen.Hash, &gen.ConfirmedSuccess)
	if err != nil {
		return nil
	}
	return &gen
}

func (mgr *DstTxManager) IsDstTxExist(srcId int64, action string, version int32) bool {
	var id int64
	err := mgr.db.QueryRow("SELECT id FROM t_dst_transaction where src_action = ? and src_id = ? and src_version = ?", action, srcId, version).Scan(&id)
	return err == nil
}

func (mgr *DstTxManager) GetDstTxConfirmGen(srcId int64, action string, version int32) int64 {
	var genId int64 = 0
	err := mgr.db.QueryRow("SELECT confirmed_gen FROM t_dst_transaction where src_action = ? and src_id = ? and src_version = ? and confirmed_gen is not null", action, srcId, version).Scan(&genId)
	if err != nil {
		return 0
	}
	return genId
}

func (mgr *DstTxManager) Save(tx *DstTx) error {
	query := `INSERT IGNORE INTO t_dst_transaction (src_action, src_id, src_version, sender, body, fee_cap, transfer_token, transfer_recipient, transfer_amount)
              VALUES (?, ?, ?, ?, ?, ?, ? , ?, ?)`

	// Execute the SQL statement with tx data
	_, err := mgr.db.Exec(query, tx.SrcAction, tx.SrcId, tx.SrcVersion, tx.Sender, tx.Body, tx.FeeCap, tx.TransferToken, tx.TransferRecipient, tx.TransferAmount)
	if err != nil {
		mgr.alerter.AlertText("failed to insert dst transaction", err)
		return err
	}

	return nil

}
