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
	OrderWeight int32
}

type ExchangeInfoManager struct {
	idExchanges   map[int32]*ExchangeInfo
	nameExchanges map[string]*ExchangeInfo
	allExchanges  []*ExchangeInfo
	db            *sql.DB
	alerter       alert.Alerter
	mutex         *sync.RWMutex
}

func NewExchangeInfoManager(db *sql.DB, alerter alert.Alerter) *ExchangeInfoManager {
	return &ExchangeInfoManager{
		idExchanges:   make(map[int32]*ExchangeInfo),
		nameExchanges: make(map[string]*ExchangeInfo),
		allExchanges:  make([]*ExchangeInfo, 0, 100),
		db:            db,
		alerter:       alerter,
		mutex:         &sync.RWMutex{},
	}
}

func (mgr *ExchangeInfoManager) GetAllExchanges() []*ExchangeInfo {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	infos := make([]*ExchangeInfo, len(mgr.allExchanges))
	copy(infos, mgr.allExchanges)
	return infos
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
	rows, err := mgr.db.Query("SELECT id, name, icon, disabled, official_url, order_weight FROM t_exchange_info")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_exchange_info error", err)
		return
	}

	defer rows.Close()

	idExchanges := make(map[int32]*ExchangeInfo)
	nameExchanges := make(map[string]*ExchangeInfo)
	allExchanges := make([]*ExchangeInfo, 0, 100)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var xchg ExchangeInfo
		if err := rows.Scan(&xchg.Id, &xchg.Name, &xchg.Icon, &xchg.Disabled, &xchg.OfficialUrl, &xchg.OrderWeight); err != nil {
			mgr.alerter.AlertText("scan t_exchange_info row error", err)
		} else {
			xchg.Name = strings.TrimSpace(xchg.Name)
			idExchanges[xchg.Id] = &xchg
			nameExchanges[strings.ToLower(xchg.Name)] = &xchg
			allExchanges = append(allExchanges, &xchg)
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
	mgr.allExchanges = allExchanges
	mgr.mutex.Unlock()
	log.Println("load all exchanges : ", counter)

}
