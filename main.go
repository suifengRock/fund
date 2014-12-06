package main

import (
	"fmt"
	"github.com/suifengRock/fund/db"
	"github.com/suifengRock/fund/spiders"
)

func main() {

	orm, err := db.mysqlEngine()
	if err != nil {
		fmt.Println(err)
		panic(err)
		return
	}
	defer orm.Close()
	err = db.sync(orm)
	if err != nil {
		panic(err)
		fmt.Println(err)
		return
	}
	baseData, _ := spiders.GetBaseFundInfoByTT()
	orm.Insert(baseData)
	for _, obj := range baseData {
		worthList, _ := spiders.GetFundHistory(obj.Code)
		orm.Insert(worthList)
	}
	fmt.Println("update......ok")

}
