package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type FundInfo struct {
	Code  uint32 `xorm:"pk"`
	Name  string
	Style string
}

type Worth struct {
	Id                int64
	Date              string //time.time
	RealWorth         float32
	AccumulativeWorth float32
	FundCode          uint32 `xorm:"index"`
}

func NewFundInfo(code uint32, name string, style string) (obj *FundInfo) {
	obj = new(FundInfo)
	obj.Code = code
	obj.Name = name
	obj.Style = style
	return
}

func NewWorth(date string, realWorth float32, accWorth float32, code uint32) (obj *Worth) {
	obj = new(Worth)
	obj.Date = date
	obj.RealWorth = realWorth
	obj.AccumulativeWorth = accWorth
	obj.FundCode = code
	return
}

func mysqlEngine() (*xorm.Engine, error) {

	conn := "test:1234@/myData?charset=utf8"

	return xorm.NewEngine("mysql", conn)
}

func sync(engine *xorm.Engine) error {
	return engine.Sync(&FundInfo{}, &Worth{})
}

func main() {

	orm, err := mysqlEngine()
	defer orm.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
		return
	}
	err = sync(orm)
	if err != nil {
		panic(err)
		fmt.Println(err)
		return
	}
	fmt.Println("...")
	resp, err := http.Get("http://www.howbuy.com/fund/fundranking/")

	if err != nil {
		fmt.Println("error: ", err)
	}
	// b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("11111")
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("22222")
	doc.Find(".result_list table tbody tr").Each(func(i int, s *goquery.Selection) {

		if i > 2 {
			fmt.Println(i)
			return
		}
		s.Find("td").Each(func(k int, sel *goquery.Selection) {
			text := sel.Text()
			fmt.Println(text)
		})

	})

}
