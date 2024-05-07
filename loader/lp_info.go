package loader

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type LpInfo struct {
	Version           int32
	TokenName         string
	FromChainName     string
	ToChainName       string
	MinValue          float64
	MaxValue          float64
	BridgeFeeRatio    float64
	MinValueStr       string
	MaxValueStr       string
	BridgeFeeRatioStr string
	MakerAddress      string
	IsDisabled        int32
}

type LpInfoManager struct {
	lpInfos    map[int32]map[string]map[string]map[string]map[string]*LpInfo
	allLpInfos []*LpInfo
	db         *sql.DB
	alerter    alert.Alerter
	mutex      *sync.RWMutex
}

func NewLpInfoManager(db *sql.DB, alerter alert.Alerter) *LpInfoManager {
	return &LpInfoManager{
		lpInfos:    make(map[int32]map[string]map[string]map[string]map[string]*LpInfo),
		allLpInfos: make([]*LpInfo, 0, 100),
		db:         db,
		alerter:    alerter,
		mutex:      &sync.RWMutex{},
	}
}

func (mgr *LpInfoManager) GetAllLpInfos() []*LpInfo {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	return mgr.allLpInfos
}

func (mgr *LpInfoManager) GetLpInfos(version int32, token string, from string, to string) (map[string]*LpInfo, bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	versionInfos, ok := mgr.lpInfos[version]
	if ok {
		ftInfos, ok := versionInfos[strings.ToLower(strings.TrimSpace(token))]
		if ok {
			infos, ok := ftInfos[strings.ToLower(strings.TrimSpace(from))]
			if ok {
				info, ok := infos[strings.ToLower(strings.TrimSpace(to))]
				return info, ok
			}

		}
	}
	return nil, false
}

func (mgr *LpInfoManager) GetLpInfo(version int32, token string, from string, to string, maker string) (*LpInfo, bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	versionInfos, ok := mgr.lpInfos[version]
	if ok {
		ftInfos, ok := versionInfos[strings.ToLower(strings.TrimSpace(token))]
		if ok {
			infos, ok := ftInfos[strings.ToLower(strings.TrimSpace(from))]
			if ok {
				info, ok := infos[strings.ToLower(strings.TrimSpace(to))]
				if ok {
					maker, ok := info[strings.ToLower(strings.TrimSpace(maker))]
					return maker, ok
				}
			}

		}
	}
	return nil, false
}

func (mgr *LpInfoManager) LoadAllLpInfo() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT version, token_name, from_chain, to_chain, maker_address, min_value, max_value, is_disabled, bridge_fee_ratio FROM t_lp_info")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_lp_info error", err)
		return
	}

	defer rows.Close()

	lpInfos := make(map[int32]map[string]map[string]map[string]map[string]*LpInfo)
	allLpInfos := make([]*LpInfo, 0, 100)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var info LpInfo
		if err := rows.Scan(&info.Version, &info.TokenName, &info.FromChainName, &info.ToChainName, &info.MakerAddress, &info.MinValueStr, &info.MaxValueStr, &info.IsDisabled, &info.BridgeFeeRatioStr); err != nil {
			mgr.alerter.AlertText("scan t_lp_info row error", err)
		} else {

			info.FromChainName = strings.TrimSpace(info.FromChainName)
			info.ToChainName = strings.TrimSpace(info.ToChainName)
			info.TokenName = strings.TrimSpace(info.TokenName)
			info.MakerAddress = strings.TrimSpace(info.MakerAddress)

			min, err := strconv.ParseFloat(info.MinValueStr, 64)
			if err != nil {
				mgr.alerter.AlertText("t_lp_info min not float", err)
				continue
			}
			max, err := strconv.ParseFloat(info.MaxValueStr, 64)
			if err != nil {
				mgr.alerter.AlertText("t_lp_info max not float", err)
				continue
			}
			bdgfee, err := strconv.ParseFloat(info.BridgeFeeRatioStr, 64)
			if err != nil {
				mgr.alerter.AlertText("t_lp_info bridge fee not float", err)
				continue
			}

			info.MinValue = min
			info.MaxValue = max
			info.BridgeFeeRatio = bdgfee

			versions, ok := lpInfos[info.Version]
			if !ok {
				versions = make(map[string]map[string]map[string]map[string]*LpInfo)
				lpInfos[info.Version] = versions
			}

			ftInfos, ok := versions[strings.ToLower(info.TokenName)]
			if !ok {
				ftInfos = make(map[string]map[string]map[string]*LpInfo)
				versions[strings.ToLower(info.TokenName)] = ftInfos
			}
			infos, ok := ftInfos[strings.ToLower(info.FromChainName)]
			if !ok {
				infos = make(map[string]map[string]*LpInfo)
				ftInfos[strings.ToLower(info.FromChainName)] = infos
			}
			makers, ok := infos[strings.ToLower(info.ToChainName)]
			if !ok {
				makers = make(map[string]*LpInfo)
				infos[strings.ToLower(info.ToChainName)] = makers
			}
			makers[strings.ToLower(info.MakerAddress)] = &info
			allLpInfos = append(allLpInfos, &info)
			counter++
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_lp_info row error", err)
		return
	}

	mgr.mutex.Lock()
	mgr.lpInfos = lpInfos
	mgr.allLpInfos = allLpInfos
	mgr.mutex.Unlock()
	log.Println("load all lp info: ", counter)

}
