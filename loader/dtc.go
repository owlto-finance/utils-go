package loader

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type Dtc struct {
	TokenName     string
	FromChainName string
	ToChainName   string
	DtcLv1        float64
	DtcLv2        float64
	DtcLv3        float64
	DtcLv4        float64
	AmountLv1     float64
	AmountLv2     float64
	AmountLv3     float64
	AmountLv4     float64
}

type DtcManager struct {
	tokenFromToDtcs map[string]map[string]map[string]*Dtc

	db      *sql.DB
	alerter alert.Alerter
	mutex   *sync.RWMutex
}

func NewDtcManager(db *sql.DB, alerter alert.Alerter) *DtcManager {
	return &DtcManager{
		tokenFromToDtcs: make(map[string]map[string]map[string]*Dtc),

		db:      db,
		alerter: alerter,
		mutex:   &sync.RWMutex{},
	}
}

func (mgr *DtcManager) GetDtc(token string, from string, to string) (*Dtc, bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	ftInfos, ok := mgr.tokenFromToDtcs[strings.ToLower(strings.TrimSpace(token))]
	if ok {
		infos, ok := ftInfos[strings.ToLower(strings.TrimSpace(from))]
		if ok {
			info, ok := infos[strings.ToLower(strings.TrimSpace(to))]
			return info, ok
		}

	}
	return nil, false
}

func (mgr *DtcManager) LoadAllDtc() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT token_name, from_chain, to_chain, dtc_lv1, dtc_lv2, dtc_lv3, dtc_lv4, amount_lv1, amount_lv2, amount_lv3, amount_lv4 FROM t_dynamic_dtc")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_dynamic_dtc error", err)
		return
	}

	defer rows.Close()

	tokenFromToDtcs := make(map[string]map[string]map[string]*Dtc)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var dtc Dtc
		var sdtc1 string
		var sdtc2 string
		var sdtc3 string
		var sdtc4 string
		var samount1 string
		var samount2 string
		var samount3 string
		var samount4 string

		if err := rows.Scan(&dtc.TokenName, &dtc.FromChainName, &dtc.ToChainName, &sdtc1, &sdtc2, &sdtc3, &sdtc4, &samount1, &samount2, &samount3, &samount4); err != nil {
			mgr.alerter.AlertText("scan t_dynamic_dtc row error", err)
		} else {
			dtc.FromChainName = strings.TrimSpace(dtc.FromChainName)
			dtc.ToChainName = strings.TrimSpace(dtc.ToChainName)
			dtc.TokenName = strings.TrimSpace(dtc.TokenName)

			dtc1, err := strconv.ParseFloat(sdtc1, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_dtc dtc1 not float", err)
				continue
			}
			dtc2, err := strconv.ParseFloat(sdtc2, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_dtc dtc2 not float", err)
				continue
			}
			dtc3, err := strconv.ParseFloat(sdtc3, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_dtc dtc3 not float", err)
				continue
			}
			dtc4, err := strconv.ParseFloat(sdtc4, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_dtc dtc4 not float", err)
				continue
			}

			amount1, err := strconv.ParseFloat(samount1, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_dtc amount1 not float", err)
				continue
			}
			amount2, err := strconv.ParseFloat(samount2, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_dtc amount2 not float", err)
				continue
			}
			amount3, err := strconv.ParseFloat(samount3, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_dtc amount3 not float", err)
				continue
			}
			amount4, err := strconv.ParseFloat(samount4, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_dtc amount4 not float", err)
				continue
			}

			dtc.DtcLv1 = dtc1
			dtc.DtcLv2 = dtc2
			dtc.DtcLv3 = dtc3
			dtc.DtcLv4 = dtc4
			dtc.AmountLv1 = amount1
			dtc.AmountLv2 = amount2
			dtc.AmountLv3 = amount3
			dtc.AmountLv4 = amount4

			ftInfos, ok := tokenFromToDtcs[strings.ToLower(dtc.TokenName)]
			if !ok {
				ftInfos = make(map[string]map[string]*Dtc)
				tokenFromToDtcs[strings.ToLower(dtc.TokenName)] = ftInfos
			}
			infos, ok := ftInfos[strings.ToLower(dtc.FromChainName)]
			if !ok {
				infos = make(map[string]*Dtc)
				ftInfos[strings.ToLower(dtc.FromChainName)] = infos
			}
			infos[strings.ToLower(dtc.ToChainName)] = &dtc

			counter++
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_dynamic_dtc row error", err)
		return
	}

	mgr.mutex.Lock()
	mgr.tokenFromToDtcs = tokenFromToDtcs
	mgr.mutex.Unlock()
	log.Println("load all dtc: ", counter)

}

func (mgr *DtcManager) GetIncludedDtc(tokenName string, fromChainName string, toChainName string, value float64) (float64, bool) {
	dtc, ok := mgr.GetDtc(tokenName, fromChainName, toChainName)
	if !ok {
		return 0, false
	}

	var dtcValue float64 = 0
	if value < (dtc.AmountLv1 + dtc.DtcLv1) {
		dtcValue = dtc.DtcLv1
	} else if value < (dtc.AmountLv2 + dtc.DtcLv2) {
		dtcValue = dtc.DtcLv2
	} else if value < (dtc.AmountLv3 + dtc.DtcLv3) {
		dtcValue = dtc.DtcLv3
	} else {
		dtcValue = dtc.DtcLv4
	}
	return dtcValue, true
}
