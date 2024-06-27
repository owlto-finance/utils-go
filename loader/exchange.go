package loader

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type ExchangeInfo struct {
	Id          int32
	Name        string
	Icon        string
	Disabled    int8
	OfficialUrl string
}

type ExchangeInfoManager struct {
	idExchanges   map[int32]*ExchangeInfo
	nameExchanges map[string]*ExchangeInfo
	db            *sql.DB
	alerter       alert.Alerter
	mutex         *sync.RWMutex
}

func NewExchangeInfoManager(db *sql.DB, alerter alert.Alerter) *ExchangeInfoManager {
	return &ExchangeInfoManager{
		idExchanges:   make(map[int32]*ExchangeInfo),
		nameExchanges: make(map[string]*ExchangeInfo),
		db:            db,
		alerter:       alerter,
		mutex:         &sync.RWMutex{},
	}
}

func (mgr *ExchangeInfoManager) GetExchangeInfoById(id int32) (*ExchangeInfo, bool) {
	mgr.mutex.RLock()
	xchg, ok := mgr.idExchanges[id]
	mgr.mutex.RUnlock()
	return xchg, ok
}
func (mgr *ExchangeInfoManager) GetExchangeInfoByName(name string) (*ExchangeInfo, bool) {
	mgr.mutex.RLock()
	xchg, ok := mgr.nameExchanges[strings.ToLower(strings.TrimSpace(name))]
	mgr.mutex.RUnlock()
	return xchg, ok
}

func (mgr *ExchangeInfoManager) LoadAllExchanges() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT id, name, icon, disabled, official_url FROM t_exchange_info")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_exchange_info error", err)
		return
	}

	defer rows.Close()

	idExchanges := make(map[int32]*ExchangeInfo)
	nameExchanges := make(map[string]*ExchangeInfo)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var xchg ExchangeInfo
		if err := rows.Scan(&xchg.Id, &xchg.Name, &xchg.Icon, &xchg.Disabled, &xchg.OfficialUrl); err != nil {
			mgr.alerter.AlertText("scan t_exchange_info row error", err)
		} else {
			xchg.Name = strings.TrimSpace(xchg.Name)
			idExchanges[xchg.Id] = &xchg
			nameExchanges[strings.ToLower(xchg.Name)] = &xchg
			counter++
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_exchange_info row error", err)
		return
	}

	mgr.mutex.Lock()
	mgr.idExchanges = idExchanges
	mgr.nameExchanges = nameExchanges
	mgr.mutex.Unlock()
	log.Println("load all exchanges : ", counter)

}
