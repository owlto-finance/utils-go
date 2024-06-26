package loader

import (
	"database/sql"
	"log"
	"sort"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type ChannelCommissionRatio struct {
	txCount int64
	ratio   int64
}

type ChannelCommissionRatioManager struct {
	channelidToCountToRatio map[int64]map[int64]int64
	channelidToRatioArr     map[int64][]ChannelCommissionRatio

	db      *sql.DB
	alerter alert.Alerter
	mutex   *sync.RWMutex
}

func NewChannelCommissionRatioManager(db *sql.DB, alerter alert.Alerter) *ChannelCommissionRatioManager {
	return &ChannelCommissionRatioManager{
		channelidToCountToRatio: make(map[int64]map[int64]int64),
		channelidToRatioArr:     make(map[int64][]ChannelCommissionRatio),

		db:      db,
		alerter: alerter,
		mutex:   &sync.RWMutex{},
	}
}

func (mgr *ChannelCommissionRatioManager) GetRatioByChannelidAndCount(channelid int64, txcount int64) (int64, bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	if sortList, isexist := mgr.channelidToRatioArr[channelid]; !isexist {
		return 0, false
	} else {
		for _, kv := range sortList {
			if txcount < kv.txCount {
				return kv.ratio, true
			}
		}
		return sortList[len(sortList)-1].ratio, true
	}
}

func (mgr *ChannelCommissionRatioManager) LoadAllCommissionRatio() {
	rows, err := mgr.db.Query("select channel_id, tx_count, commission_ratio from t_channel_commission_ratio order by tx_count asc")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_channel_commission_ratio error", err)
		return
	}

	defer rows.Close()

	channelidToCountToRatio := make(map[int64]map[int64]int64)
	channelidToRatioArr := make(map[int64][]ChannelCommissionRatio)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var channelID, txCount, ratio int64

		if err := rows.Scan(&channelID, &txCount, &ratio); err != nil {
			mgr.alerter.AlertText("scan t_channel_commission_ratio row error", err)
		} else {
			channelidToCountToRatio[channelID] = make(map[int64]int64)
			channelidToCountToRatio[channelID][txCount] = ratio

			ratioArr, exist := channelidToRatioArr[channelID]
			if !exist {
				ratioArr = make([]ChannelCommissionRatio, 0)
			}
			ratioArr = append(ratioArr, ChannelCommissionRatio{txCount, ratio})
			channelidToRatioArr[channelID] = ratioArr

			counter++
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_channel_commission_ratio row error", err)
		return
	}

	for k, _ := range channelidToRatioArr {
		sort.Slice(channelidToRatioArr[k], func(i, j int) bool {
			return channelidToRatioArr[k][i].txCount < channelidToRatioArr[k][j].txCount
		})
	}

	mgr.mutex.Lock()
	mgr.channelidToCountToRatio = channelidToCountToRatio
	mgr.channelidToRatioArr = channelidToRatioArr
	mgr.mutex.Unlock()
	log.Println("load all channel commission ratio: ", counter)

}
