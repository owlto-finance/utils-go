package loader

import (
	"database/sql"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
	"github.com/owlto-finance/utils-go/util"
)

type BridgeFee struct {
	TokenName         string
	FromChainName     string
	ToChainName       string
	BridgeFeeRatioLv1 int64
	BridgeFeeRatioLv2 int64
	BridgeFeeRatioLv3 int64
	BridgeFeeRatioLv4 int64
	AmountLv1         float64
	AmountLv2         float64
	AmountLv3         float64
	AmountLv4         float64
	KeepDecimal       int32

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
				mgr.alerter.AlertText("t_dynamic_bridge_fee keep decimal not found: token "+bridgeFee.TokenName+" chain "+bridgeFee.FromChainName, err)
				continue
			} else if kdexist {
				bridgeFee.KeepDecimal = int32(dbKeepDecimal)
			} else if ok {
				bridgeFee.KeepDecimal = int32(tokenInfo.Decimals)
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

func (mgr *BridgeFeeManager) FromUiString(amount *big.Int, bridgeFee string, decimal int32, keepDecimal int32) *big.Int {
	value := big.NewInt(0)
	value.Add(value, amount)

	if bridgeFee != "" {
		bridgeFeeValue, err := util.FromUiString(bridgeFee, decimal)
		if err == nil {
			bridgeFeeAmount := new(big.Int)
			bridgeFeeAmount.Mul(value, bridgeFeeValue)
			bridgeFeeAmount.Div(value, big.NewInt(100000000))

			scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal-keepDecimal)), nil)
			modAmount := new(big.Int)
			modAmount.Mod(bridgeFeeAmount, scale)
			bridgeFeeAmount.Sub(bridgeFeeAmount, modAmount)

			value.Sub(value, bridgeFeeAmount)
		}
	}

	return value
}

func (mgr *BridgeFeeManager) GetIncludedBridgeFeeBigInt(tokenName string, fromChainName string, toChainName string, value *big.Int, decimal int32) (int64, bool) {
	bridgeFee, ok := mgr.GetBridgeFee(tokenName, fromChainName, toChainName)
	if !ok {
		return 0, false
	}

	keepDecimal := decimal
	if bridgeFee.KeepDecimal < decimal {
		keepDecimal = bridgeFee.KeepDecimal
	}

	AmountLv1BigInt, err := util.FromUiString(bridgeFee.AmountLv1Str, decimal)
	if err != nil {
		return 0, false
	}
	AmountLv2BigInt, err := util.FromUiString(bridgeFee.AmountLv2Str, decimal)
	if err != nil {
		return 0, false
	}
	AmountLv3BigInt, err := util.FromUiString(bridgeFee.AmountLv3Str, decimal)
	if err != nil {
		return 0, false
	}

	if AmountLv1BigInt.Cmp(mgr.FromUiString(value, strconv.FormatInt(bridgeFee.BridgeFeeRatioLv1, 10), decimal, keepDecimal)) > 0 {
		return bridgeFee.BridgeFeeRatioLv1, true
	} else if AmountLv2BigInt.Cmp(mgr.FromUiString(value, strconv.FormatInt(bridgeFee.BridgeFeeRatioLv2, 10), decimal, keepDecimal)) > 0 {
		return bridgeFee.BridgeFeeRatioLv2, true
	} else if AmountLv3BigInt.Cmp(mgr.FromUiString(value, strconv.FormatInt(bridgeFee.BridgeFeeRatioLv3, 10), decimal, keepDecimal)) > 0 {
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
