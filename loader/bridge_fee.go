package loader

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type BridgeFee struct {
	TokenName     string
	FromChainName string
	ToChainName   string
	Bridge_fee_ratio_lv1        int64 
	Bridge_fee_ratio_lv2        int64
	Bridge_fee_ratio_lv3        int64
	Bridge_fee_ratio_lv4        int64
	AmountLv1     float64
	AmountLv2     float64
	AmountLv3     float64
	AmountLv4     float64

	AmountLv1Str string
	AmountLv2Str string
	AmountLv3Str string
	AmountLv4Str string
}

type BridgeFeeManager struct {
	tokenFromToBridgeFees map[string]map[string]map[string]*BridgeFee

	db      *sql.DB
	alerter alert.Alerter
	mutex   *sync.RWMutex
}

func NewBridgeFeeManager(db *sql.DB, alerter alert.Alerter) *BridgeFeeManager {
	return &BridgeFeeManager{
    tokenFromToBridgeFees: make(map[string]map[string]map[string]*BridgeFee),

		db:      db,
		alerter: alerter,
		mutex:   &sync.RWMutex{},
	}
}

func (mgr *BridgeFeeManager) GetBridgeFee(token string, from string, to string) (*BridgeFee, bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	ftInfos, ok := mgr.tokenFromToBridgeFees[strings.ToLower(strings.TrimSpace(token))]
	if ok {
		infos, ok := ftInfos[strings.ToLower(strings.TrimSpace(from))]
		if ok {
			info, ok := infos[strings.ToLower(strings.TrimSpace(to))]
			return info, ok
		}

	}
	return nil, false
}

func (mgr *BridgeFeeManager) LoadAllBridgeFee() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT token_name, from_chain, to_chain, bridge_fee_ratio_lv1, bridge_fee_ratio_lv2, bridge_fee_ratio_lv3, bridge_fee_ratio_lv4, amount_lv1, amount_lv2, amount_lv3, amount_lv4 FROM t_dynamic_bridge_fee")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_dynamic_bridge_fee error", err)
		return
	}

	defer rows.Close()

	tokenFromToBridgeFees := make(map[string]map[string]map[string]*BridgeFee)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var bridgeFee BridgeFee

		if err := rows.Scan(&bridgeFee.TokenName, &bridgeFee.FromChainName, &bridgeFee.ToChainName, &bridgeFee.Bridge_fee_ratio_lv1, &bridgeFee.Bridge_fee_ratio_lv2, &bridgeFee.Bridge_fee_ratio_lv3, &bridgeFee.Bridge_fee_ratio_lv4, &bridgeFee.AmountLv1Str, &bridgeFee.AmountLv2Str, &bridgeFee.AmountLv3Str, &bridgeFee.AmountLv4Str); err != nil {
			mgr.alerter.AlertText("scan t_dynamic_bridge_fee row error", err)
		} else {
			bridgeFee.FromChainName = strings.TrimSpace(bridgeFee.FromChainName)
			bridgeFee.ToChainName = strings.TrimSpace(bridgeFee.ToChainName)
			bridgeFee.TokenName = strings.TrimSpace(bridgeFee.TokenName)

			amount1, err := strconv.ParseFloat(bridgeFee.AmountLv1Str, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_bridge_fee amount1 not float", err)
				continue
			}
			amount2, err := strconv.ParseFloat(bridgeFee.AmountLv2Str, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_bridge_fee amount2 not float", err)
				continue
			}
			amount3, err := strconv.ParseFloat(bridgeFee.AmountLv3Str, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_bridge_fee amount3 not float", err)
				continue
			}
			amount4, err := strconv.ParseFloat(bridgeFee.AmountLv4Str, 64)
			if err != nil {
				mgr.alerter.AlertText("t_dynamic_bridge_fee amount4 not float", err)
				continue
			}

			bridgeFee.AmountLv1 = amount1
			bridgeFee.AmountLv2 = amount2
			bridgeFee.AmountLv3 = amount3
			bridgeFee.AmountLv4 = amount4

			ftInfos, ok := tokenFromToBridgeFees[strings.ToLower(bridgeFee.TokenName)]
			if !ok {
				ftInfos = make(map[string]map[string]*BridgeFee)
				tokenFromToBridgeFees[strings.ToLower(bridgeFee.TokenName)] = ftInfos
			}
			infos, ok := ftInfos[strings.ToLower(bridgeFee.FromChainName)]
			if !ok {
				infos = make(map[string]*BridgeFee)
				ftInfos[strings.ToLower(bridgeFee.FromChainName)] = infos
			}
			infos[strings.ToLower(bridgeFee.ToChainName)] = &bridgeFee

			counter++
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_dynamic_bridge_fee row error", err)
		return
	}

	mgr.mutex.Lock()
	mgr.tokenFromToBridgeFees = tokenFromToBridgeFees
	mgr.mutex.Unlock()
	log.Println("load all bridge fee: ", counter)

}

func (mgr *BridgeFeeManager) GetIncludedBridgeFee(tokenName string, fromChainName string, toChainName string, value float64, dtc float64) (int64, bool) {
	bridgeFee, ok := mgr.GetBridgeFee(tokenName, fromChainName, toChainName)
	if !ok {
		return 0, false
	}
  //parseFloat(bridgefeeItem.amount_lv1) > (parseFloat(amount) - parseFloat(dtc)) * (1 - parseFloat(bridgefeeItem.bridge_fee_ratio_lv1.toString())/1000000.0*0.01)
  if bridgeFee.AmountLv1 > (value - dtc) * (1 - float64(bridgeFee.Bridge_fee_ratio_lv1) / 1000000 * 0.01) {
    return bridgeFee.Bridge_fee_ratio_lv1, true
  } else if bridgeFee.AmountLv2 > (value - dtc) * (1 - float64(bridgeFee.Bridge_fee_ratio_lv2) / 1000000 * 0.01) {
    return bridgeFee.Bridge_fee_ratio_lv2, true
  } else if bridgeFee.AmountLv3 > (value - dtc) * (1 - float64(bridgeFee.Bridge_fee_ratio_lv3) / 1000000 * 0.01) {
    return bridgeFee.Bridge_fee_ratio_lv3, true
  } else {
    return bridgeFee.Bridge_fee_ratio_lv4, true
  }
}

func (mgr *BridgeFeeManager) GetBridgeFeeNotIncluded(tokenName string, fromChainName string, toChainName string, value float64) (int64, bool) {
	bridgeFee, ok := mgr.GetBridgeFee(tokenName, fromChainName, toChainName)
	if !ok {
		return 0, false
	}

	if value < bridgeFee.AmountLv1 {
		return bridgeFee.Bridge_fee_ratio_lv1, true
	} else if value < bridgeFee.AmountLv2 {
		return bridgeFee.Bridge_fee_ratio_lv2, true
	} else if value < bridgeFee.AmountLv3 {
		return bridgeFee.Bridge_fee_ratio_lv3, true
	} else {
		return bridgeFee.Bridge_fee_ratio_lv4, true
	}
}
