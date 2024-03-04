package loader

import (
	"database/sql"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/owlto-finance/utils-go/alert"
)

type ChainInfo struct {
	Id          int64
	ChainId     string
	Name        string
	IsTestnet   int8
	RpcEndPoint string
	Disabled    int8
	Client      *ethclient.Client
}

type ChainInfoManager struct {
	idChains      map[int64]*ChainInfo
	chainIdChains map[string]*ChainInfo
	nameChains    map[string]*ChainInfo
	db            *sql.DB
	alerter       alert.Alerter
	mutex         *sync.RWMutex
}

func NewChainInfoManager(db *sql.DB, alerter alert.Alerter) ChainInfoManager {
	return ChainInfoManager{
		idChains:      make(map[int64]*ChainInfo),
		chainIdChains: make(map[string]*ChainInfo),
		nameChains:    make(map[string]*ChainInfo),
		db:            db,
		alerter:       alerter,
		mutex:         &sync.RWMutex{},
	}
}

func (mgr *ChainInfoManager) GetChainInfoById(id int64) (*ChainInfo, bool) {
	mgr.mutex.RLock()
	chain, ok := mgr.idChains[id]
	mgr.mutex.RUnlock()
	return chain, ok
}
func (mgr *ChainInfoManager) GetChainInfoByChainId(chainId string) (*ChainInfo, bool) {
	mgr.mutex.RLock()
	chain, ok := mgr.chainIdChains[chainId]
	mgr.mutex.RUnlock()
	return chain, ok
}
func (mgr *ChainInfoManager) GetChainInfoByName(name string) (*ChainInfo, bool) {
	mgr.mutex.RLock()
	chain, ok := mgr.nameChains[name]
	mgr.mutex.RUnlock()
	return chain, ok
}

func (mgr *ChainInfoManager) LoadAllChains() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT id, chainid, name, is_testnet,rpc_end_point, disabled FROM t_chain_info")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_chain_info error", err)
		return
	}

	defer rows.Close()
	idChains := make(map[int64]*ChainInfo)
	chainIdChains := make(map[string]*ChainInfo)
	nameChains := make(map[string]*ChainInfo)

	counter := 0
	// Iterate over the result set
	for rows.Next() {
		var chain ChainInfo

		if err := rows.Scan(&chain.Id, &chain.ChainId, &chain.Name, &chain.IsTestnet, &chain.RpcEndPoint, &chain.Disabled); err != nil {
			mgr.alerter.AlertText("scan t_chain_info row error", err)
		} else {
			chain.Client, err = ethclient.Dial(chain.RpcEndPoint)
			if err != nil {
				mgr.alerter.AlertText("create client error", err)
				continue
			}
			idChains[chain.Id] = &chain
			chainIdChains[chain.ChainId] = &chain
			nameChains[chain.Name] = &chain
			counter++
		}
	}

	mgr.mutex.Lock()
	mgr.idChains = idChains
	mgr.chainIdChains = chainIdChains
	mgr.nameChains = nameChains
	mgr.mutex.Unlock()

	log.Println("load all chain info: ", counter)
	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_chain_info row error", err)
		return
	}

}
