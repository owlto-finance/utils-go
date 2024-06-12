package loader

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type PopularList struct {
	ChainName     string
	PopularWeight       int32
}

type PopularListManager struct {
	chainToPopularList map[string]*PopularList

	db      *sql.DB
	alerter alert.Alerter
	mutex   *sync.RWMutex
}

func NewPopularListManager(db *sql.DB, alerter alert.Alerter) *PopularListManager {
	return &PopularListManager{
		chainToPopularList: make(map[string]*PopularList),

		db:      db,
		alerter: alerter,
		mutex:   &sync.RWMutex{},
	}
}

func (mgr *PopularListManager) GetPopularWeight(chain string) (int32, bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
  plInfo, ok := mgr.chainToPopularList[strings.ToLower(strings.TrimSpace(chain))]
	if ok {
		return plInfo.PopularWeight, ok
	}
	return -1, false
}

func (mgr *PopularListManager) LoadAllPopularList() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT chain_name, popular_weight FROM t_popular_list")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_popular_list error", err)
		return
	}

	defer rows.Close()

  chainToPopularList := make(map[string]*PopularList)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var popularList PopularList

		if err := rows.Scan(&popularList.ChainName, &popularList.PopularWeight); err != nil {
			mgr.alerter.AlertText("scan t_popular_list row error", err)
		} else {
      popularList.ChainName = strings.TrimSpace(popularList.ChainName)

			chainToPopularList[strings.ToLower(popularList.ChainName)] = &popularList

			counter++
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_popular_list row error", err)
		return
	}

	mgr.mutex.Lock()
  mgr.chainToPopularList = chainToPopularList
	mgr.mutex.Unlock()
	log.Println("load all popular list: ", counter)

}

