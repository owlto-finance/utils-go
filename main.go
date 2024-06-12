package main

import (
		"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
		"github.com/owlto-finance/utils-go/alert"
	"github.com/owlto-finance/utils-go/util"
  "github.com/owlto-finance/utils-go/loader"
)

func main() {

		db, kde := sql.Open("mysql", "root:SDf#78@kJ9@tcp(10.111.64.3:3306)/db_paypto")
		if kde != nil {
			fmt.Println(kde)
			return
		}
		alerter := alert.NewLarkAlerter("https://open.larksuite.com/open-apis/bot/v2/hook/9a215f4f-b379-411d-99a3-2fb9a17552fa")
		popularListMgr := loader.NewPopularListManager(db, alerter)
    popularListMgr.LoadAllPopularList()
	//	chainInfoMgr.LoadAllChains()
	//	tokenInfoMgr := loader.NewTokenInfoManager(db, alerter)
	//	tokenInfoMgr.MergeNativeTokens(*chainInfoMgr)
	//	tokenInfoMgr.LoadAllToken()
	//	_ = loader.NewBridgeFeeManager(kddb, alerter)

	amountBigInt, err := util.FromUiString("123", 2)
	if err != nil {
		fmt.Println("???")
		return
	}
	amount2BigInt, err := util.FromUiString("321", 2)
	amount3BigInt, err := util.FromUiString("1", 2)
	amountBigInt.Add(amountBigInt, amount2BigInt)
	amountBigInt.Div(amountBigInt, amount3BigInt)

	fmt.Println(amountBigInt)

}
