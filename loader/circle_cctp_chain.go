package loader

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type CircleCctpChain struct {
	ChainId            int32
	MinValue           string
	Domain             int32
	TokenMessenger     string
	MessageTransmitter string
}

func (ccc *CircleCctpChain) GetMinValueUnit() *big.Int {
	min, _ := new(big.Int).SetString(ccc.MinValue, 0)
	result := new(big.Int).Mul(min, big.NewInt(1000000))
	return result
}

type CircleCctpChainManager struct {
	chainIdChains map[int32]*CircleCctpChain
	db            *sql.DB
	alerter       alert.Alerter
	mutex         *sync.RWMutex
}

func NewCircleCctpChainManager(db *sql.DB, alerter alert.Alerter) *CircleCctpChainManager {
	return &CircleCctpChainManager{
		chainIdChains: make(map[int32]*CircleCctpChain),
		db:            db,
		alerter:       alerter,
		mutex:         &sync.RWMutex{},
	}
}

func (mgr *CircleCctpChainManager) GetDtcUnit(srcChainId int32, dstChainId int32) *big.Int {
	if srcChainId == 1 || dstChainId == 1 {
		return big.NewInt(50000000)
	}
	return big.NewInt(5000000)
}

func (mgr *CircleCctpChainManager) GetChainByChainId(id int32) (*CircleCctpChain, bool) {
	mgr.mutex.RLock()
	chain, ok := mgr.chainIdChains[id]
	mgr.mutex.RUnlock()
	return chain, ok
}

func (mgr *CircleCctpChainManager) GetChainIds() []int32 {
	mgr.mutex.RLock()
	chainIds := make([]int32, 0, len(mgr.chainIdChains))

	// Iterate over the map and extract keys
	for chainId := range mgr.chainIdChains {
		chainIds = append(chainIds, chainId)
	}
	mgr.mutex.RUnlock()
	return chainIds
}

func (mgr *CircleCctpChainManager) LoadAllChains() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT chainid, min_value, domain, token_messenger, message_transmitter FROM t_cctp_support_chain")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_cctp_support_chain error", err)
		return
	}

	defer rows.Close()

	chainIdChains := make(map[int32]*CircleCctpChain)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var chain CircleCctpChain
		if err := rows.Scan(&chain.ChainId, &chain.MinValue, &chain.Domain, &chain.TokenMessenger, &chain.MessageTransmitter); err != nil {
			mgr.alerter.AlertText("scan t_cctp_support_chain row error", err)
		} else {
			chain.MessageTransmitter = strings.TrimSpace(chain.MessageTransmitter)
			chain.TokenMessenger = strings.TrimSpace(chain.TokenMessenger)
			chain.MinValue = strings.TrimSpace(chain.MinValue)
			_, ok := new(big.Int).SetString(chain.MinValue, 0)
			if !ok {
				mgr.alerter.AlertText("scan t_cctp_support_chain min value error ", fmt.Errorf("id: %d, min value: %s", chain.ChainId, chain.MinValue))
			}

			chainIdChains[chain.ChainId] = &chain
			counter++
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_cctp_support_chain row error", err)
		return
	}

	mgr.mutex.Lock()
	mgr.chainIdChains = chainIdChains
	mgr.mutex.Unlock()
	log.Println("load all cctp chain: ", counter)

}
