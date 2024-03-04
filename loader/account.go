package loader

import (
	"database/sql"
	"log"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type Account struct {
	Id          int64
	ChainInfoId int64
	Address     string
}

type AccountManager struct {
	idAccounts          map[int64]*Account
	chainInfoIdAccounts map[int64]*Account
	db                  *sql.DB
	alerter             alert.Alerter
	mutex               *sync.RWMutex
}

func NewAccountManager(db *sql.DB, alerter alert.Alerter) *AccountManager {
	return &AccountManager{
		idAccounts:          make(map[int64]*Account),
		chainInfoIdAccounts: make(map[int64]*Account),
		db:                  db,
		alerter:             alerter,
		mutex:               &sync.RWMutex{},
	}
}

func (mgr *AccountManager) GetAccountById(id int64) (*Account, bool) {
	mgr.mutex.RLock()
	acc, ok := mgr.idAccounts[id]
	mgr.mutex.RUnlock()
	return acc, ok
}

func (mgr *AccountManager) GetAccountByChainInfoId(id int64) (*Account, bool) {
	mgr.mutex.RLock()
	acc, ok := mgr.chainInfoIdAccounts[id]
	mgr.mutex.RUnlock()
	return acc, ok
}

func (mgr *AccountManager) LoadAllAccounts() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT id, chain_id, address FROM t_account")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_account error", err)
		return
	}

	defer rows.Close()

	idAccounts := make(map[int64]*Account)
	chainInfoIdAccounts := make(map[int64]*Account)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var acc Account
		if err := rows.Scan(&acc.Id, &acc.ChainInfoId, &acc.Address); err != nil {
			mgr.alerter.AlertText("scan t_account row error", err)
		} else {
			idAccounts[acc.Id] = &acc
			chainInfoIdAccounts[acc.ChainInfoId] = &acc
			counter++
		}
	}

	mgr.mutex.Lock()
	mgr.idAccounts = idAccounts
	mgr.chainInfoIdAccounts = chainInfoIdAccounts
	mgr.mutex.Unlock()
	log.Println("load all account: ", counter)

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_account row error", err)
		return
	}

}
