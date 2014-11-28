package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/qiniu/iconv"
)

type FundInfo struct {
	Id    int64
	Code  string `xorm:"index unique"`
	Name  string `xorm:"not null"`
	Style string
}

type Worth struct {
	Id                int64
	Date              string
	RealWorth         float32
	AccumulativeWorth float32
	FundCode          string `xorm:"index"`
}

func NewFundInfo(code string, name string, style string) (obj *FundInfo) {
	obj = new(FundInfo)
	obj.Code = code
	obj.Name = name
	obj.Style = style
	return
}

func NewWorth(date string, realWorth float32, accWorth float32, code string) (obj *Worth) {
	obj = new(Worth)
	obj.Date = date
	obj.RealWorth = realWorth
	obj.AccumulativeWorth = accWorth
	obj.FundCode = code
	return
}

func mysqlEngine() (*xorm.Engine, error) {

	addr := os.Getenv("DB_PORT_3306_TCP_ADDR")
	port := os.Getenv("DB_PORT_3306_TCP_PORT")
	proto := os.Getenv("DB_PORT_3306_TCP_PROTO")
	user := os.Getenv("DB_ENV_MYSQL_USER")
	password := os.Getenv("DB_ENV_MYSQL_PASSWORD")
	database := os.Getenv("DB_ENV_MYSQL_DATABASE")

	conn := "test:1234@/myData?charset=utf8"

	if addr != "" {
		conn = fmt.Sprintf("%v:%v@%v(%v:%v)/%v?charset=utf8", user, password, proto, addr, port, database)
		fmt.Println("the connection is " + conn)
	}

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
	baseData, _ := getBaseFundInfoByTT()
	orm.Insert(baseData)
	fmt.Println("update......ok")

}

func getBaseFundInfoByTT() ([]*FundInfo, error) {

	urlFund := "http://fund.eastmoney.com/fund.html"
	doc, err := goquery.NewDocument(urlFund)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cd, err := iconv.Open("utf-8", "gbk")
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return nil, err
	}
	defer cd.Close()
	baseData := make([]*FundInfo, 0)
	doc.Find("#oTable tbody tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		code := ""
		name := ""
		s.Find("td").Each(func(key int, sel *goquery.Selection) {

			if key == 3 {
				code = sel.Text()
			}
			if key == 4 {
				name = sel.Find("nobr a").First().Text()
			}
		})
		obj := NewFundInfo(code, cd.ConvString(name), "")
		baseData = append(baseData, obj)
	})

	return baseData, nil
}
