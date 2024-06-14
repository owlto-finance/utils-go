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
	PopularWeight map[string]int32
}

type PopularListManager struct {
	chainToPopularList map[string]PopularList

	db      *sql.DB
	alerter alert.Alerter
	mutex   *sync.RWMutex
}

func NewPopularListManager(db *sql.DB, alerter alert.Alerter) *PopularListManager {
	return &PopularListManager{
		chainToPopularList: make(map[string]PopularList),

		db:      db,
		alerter: alerter,
		mutex:   &sync.RWMutex{},
	}
}

func (mgr *PopularListManager) GetPopularWeight(weights map[string]int32, chain string) bool {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	plInfo, ok := mgr.chainToPopularList[strings.ToLower(strings.TrimSpace(chain))]
	if ok {
		weights = plInfo.PopularWeight
		return ok
	}
	return false
}

func (mgr *PopularListManager) LoadAllPopularList() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT chain_name, popular_weight, tag FROM t_popular_list")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_popular_list error", err)
		return
	}

	defer rows.Close()

	chainToPopularList := make(map[string]PopularList)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var weight int32
		var chainName, tag string

		if err := rows.Scan(&chainName, &weight, &tag); err != nil {
			mgr.alerter.AlertText("scan t_popular_list row error", err)
		} else {
			chainName = strings.ToLower(strings.TrimSpace(chainName))

			var popularList PopularList
      popularList.PopularWeight = make(map[string]int32)
			if pl, ok := chainToPopularList[chainName]; ok {
				popularList = pl
			}
			popularList.ChainName = chainName
			popularList.PopularWeight[strings.TrimSpace(tag)] = weight

			chainToPopularList[chainName] = popularList

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
