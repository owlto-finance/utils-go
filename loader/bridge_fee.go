package loader

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"sync"
  "math"

	"github.com/owlto-finance/utils-go/alert"
)

type BridgeFee struct {
	TokenName     string
	FromChainName string
	ToChainName   string
	BridgeFeeRatioLv1        int64 
	BridgeFeeRatioLv2        int64
	BridgeFeeRatioLv3        int64
	BridgeFeeRatioLv4        int64
	AmountLv1     float64
	AmountLv2     float64
	AmountLv3     float64
	AmountLv4     float64
  KeepDecimal   int64

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

func (mgr *BridgeFeeManager) LoadAllBridgeFee(tokenInfoMgr TokenInfoManager) {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT token_name, from_chain, to_chain, bridge_fee_ratio_lv1, bridge_fee_ratio_lv2, bridge_fee_ratio_lv3, bridge_fee_ratio_lv4, amount_lv1, amount_lv2, amount_lv3, amount_lv4 FROM t_dynamic_bridge_fee")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_dynamic_bridge_fee error", err)
		return
	}

	kdrows, kderr := mgr.db.Query("SELECT token, keep_decimal FROM t_bridge_fee_decimal")
  if kderr != nil {
    mgr.alerter.AlertText("select t_bridge_fee_decimal error", kderr)
  }

	defer rows.Close()
  defer kdrows.Close()

  tokenDecimal := make(map[string]int64)
  for kdrows.Next() {
    var tokenName string
    var keepDecimal int64
		if err := kdrows.Scan(&tokenName, &keepDecimal); err != nil {
			mgr.alerter.AlertText("scan t_bridge_fee_decimal row error", err)
		} else {
      tokenName = strings.TrimSpace(tokenName)
      tokenDecimal[strings.ToLower(tokenName)] = keepDecimal
    }
  }

	tokenFromToBridgeFees := make(map[string]map[string]map[string]*BridgeFee)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var bridgeFee BridgeFee

		if err := rows.Scan(&bridgeFee.TokenName, &bridgeFee.FromChainName, &bridgeFee.ToChainName, &bridgeFee.BridgeFeeRatioLv1, &bridgeFee.BridgeFeeRatioLv2, &bridgeFee.BridgeFeeRatioLv3, &bridgeFee.BridgeFeeRatioLv4, &bridgeFee.AmountLv1Str, &bridgeFee.AmountLv2Str, &bridgeFee.AmountLv3Str, &bridgeFee.AmountLv4Str); err != nil {
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
      
      tokenInfo, ok := tokenInfoMgr.GetByChainNameTokenName(strings.ToLower(bridgeFee.FromChainName), strings.ToLower(bridgeFee.TokenName))
      dbKeepDecimal, kdexist := tokenDecimal[strings.ToLower(bridgeFee.TokenName)]
      if !ok && !kdexist {
        mgr.alerter.AlertText("t_dynamic_bridge_fee keep decimal not found: token " + bridgeFee.TokenName + " chain " + bridgeFee.FromChainName, err)
				continue
      } else if kdexist {
        bridgeFee.KeepDecimal = dbKeepDecimal
      } else if ok {
        bridgeFee.KeepDecimal = int64(tokenInfo.Decimals)
      }

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

func truncateFloat(f float64, prec int64) float64 {
  multiplier := math.Pow(10, float64(prec))
  return math.Floor(f * multiplier) / multiplier
}

func (mgr *BridgeFeeManager) GetIncludedBridgeFee(tokenName string, fromChainName string, toChainName string, value float64, dtc float64) (int64, bool) {
	bridgeFee, ok := mgr.GetBridgeFee(tokenName, fromChainName, toChainName)
	if !ok {
		return 0, false
	}

//  parseFloat(dynamic_bridge_fee.amount_lv1) > valueWithoutDtc - parseFloat((valueWithoutDtc * parseFloat(dynamic_bridge_fee.bridge_fee_ratio_lv1.toString())/1000000.0*0.01).toFixed(keepDecimal + 1).slice(0, -1))
  if bridgeFee.AmountLv1 > value - dtc - truncateFloat((value - dtc) * float64(bridgeFee.BridgeFeeRatioLv1) / 1000000 * 0.01, bridgeFee.KeepDecimal) {
    return bridgeFee.BridgeFeeRatioLv1, true
  } else if bridgeFee.AmountLv2 > value - dtc - truncateFloat((value - dtc) * float64(bridgeFee.BridgeFeeRatioLv2) / 1000000 * 0.01, bridgeFee.KeepDecimal) {
    return bridgeFee.BridgeFeeRatioLv2, true
  } else if bridgeFee.AmountLv3 > value - dtc - truncateFloat((value - dtc) * float64(bridgeFee.BridgeFeeRatioLv3) / 1000000 * 0.01, bridgeFee.KeepDecimal) {
    return bridgeFee.BridgeFeeRatioLv3, true
  } else {
    return bridgeFee.BridgeFeeRatioLv4, true
  }
}

func (mgr *BridgeFeeManager) GetBridgeFeeNotIncluded(tokenName string, fromChainName string, toChainName string, value float64) (int64, bool) {
	bridgeFee, ok := mgr.GetBridgeFee(tokenName, fromChainName, toChainName)
	if !ok {
		return 0, false
	}

	if value < bridgeFee.AmountLv1 {
		return bridgeFee.BridgeFeeRatioLv1, true
	} else if value < bridgeFee.AmountLv2 {
		return bridgeFee.BridgeFeeRatioLv2, true
	} else if value < bridgeFee.AmountLv3 {
		return bridgeFee.BridgeFeeRatioLv3, true
	} else {
		return bridgeFee.BridgeFeeRatioLv4, true
	}
}
