package loader

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/owlto-finance/utils-go/alert"
)

type UpdatePrice struct {
	TokenName       string
	Price           string
	UpdateTimestamp string
}

type UpdatePriceManager struct {
	tokens map[string]*UpdatePrice

	db      *sql.DB
	alerter alert.Alerter
	mutex   *sync.RWMutex
}

func NewUpdatePriceManager(db *sql.DB, alerter alert.Alerter) *UpdatePriceManager {
	return &UpdatePriceManager{
		tokens: make(map[string]*UpdatePrice),

		db:      db,
		alerter: alerter,
		mutex:   &sync.RWMutex{},
	}
}

func (mgr *UpdatePriceManager) GetUpdatePrice(tokenName string) (*UpdatePrice, bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	info, ok := mgr.tokens[strings.ToLower(strings.TrimSpace(tokenName))]
	return info, ok
}

func (mgr *UpdatePriceManager) LoadAllPrice() {
	// Query the database to select only id and name fields
	rows, err := mgr.db.Query("SELECT token, price, update_timestamp FROM t_update_price")
	if err != nil || rows == nil {
		mgr.alerter.AlertText("select t_update error", err)
		return
	}
	defer rows.Close()

	tokens := make(map[string]*UpdatePrice)
	counter := 0

	for rows.Next() {
		var prices UpdatePrice

		if err := rows.Scan(&prices.TokenName, &prices.Price, &prices.UpdateTimestamp); err != nil {
			mgr.alerter.AlertText("scan t_update_price row error", err)
		} else {
			prices.TokenName = strings.TrimSpace(prices.TokenName)
			prices.Price = strings.TrimSpace(prices.Price)
			prices.UpdateTimestamp = strings.TrimSpace(prices.UpdateTimestamp)
			tokens[strings.ToLower(prices.TokenName)] = &prices
			counter++
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		mgr.alerter.AlertText("get next t_update_price row error", err)
		return
	}

	mgr.mutex.Lock()
	mgr.tokens = tokens
	mgr.mutex.Unlock()
	log.Println("load all update price: ", counter)
}
