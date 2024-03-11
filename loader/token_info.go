package loader

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type TokenInfo struct {
	TokenName    string
	ChainName    string
	TokenAddress string
	Decimals     int32
}

type TokenInfoManager struct {
	chainNameTokenAddrs map[string]map[string]*TokenInfo

	db      *sql.DB
	alerter alert.Alerter
	mutex   *sync.RWMutex
}

func NewTokenInfoManager(db *sql.DB, alerter alert.Alerter) *TokenInfoManager {
	return &TokenInfoManager{
		chainNameTokenAddrs: make(map[string]map[string]*TokenInfo),

		db:      db,
		alerter: alerter,
		mutex:   &sync.RWMutex{},
	}
}

func (mgr *TokenInfoManager) GetByChainNameTokenAddr(chainName string, tokenAddr string) (*TokenInfo, bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	tokenAddrs, ok := mgr.chainNameTokenAddrs[strings.ToLower(strings.TrimSpace(chainName))]
	if ok {
		token, ok := tokenAddrs[strings.ToLower(strings.TrimSpace(tokenAddr))]
		return token, ok
	}
	return nil, false
}

func (mgr *TokenInfoManager) LoadAllToken() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT token_name, chain_name, token_address, decimals FROM t_token_info")

	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_token_info error", err)
		return
	}

	defer rows.Close()

	chainNameTokenAddrs := make(map[string]map[string]*TokenInfo)
	counter := 0

	// Iterate over the result set
	for rows.Next() {
		var token TokenInfo

		if err := rows.Scan(&token.TokenName, &token.ChainName, &token.TokenAddress, &token.Decimals); err != nil {
			mgr.alerter.AlertText("scan t_token_info row error", err)
		} else {
			token.ChainName = strings.TrimSpace(token.ChainName)
			token.TokenAddress = strings.TrimSpace(token.TokenAddress)
			token.TokenName = strings.TrimSpace(token.TokenName)

			tokenAddrs, ok := chainNameTokenAddrs[strings.ToLower(token.ChainName)]
			if !ok {
				tokenAddrs = make(map[string]*TokenInfo)
				chainNameTokenAddrs[strings.ToLower(token.ChainName)] = tokenAddrs
			}
			tokenAddrs[strings.ToLower(token.TokenAddress)] = &token

			counter++
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_token_info row error", err)
		return
	}

	mgr.mutex.Lock()
	mgr.chainNameTokenAddrs = chainNameTokenAddrs
	mgr.mutex.Unlock()
	log.Println("load all token info: ", counter)

}
