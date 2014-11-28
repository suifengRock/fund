package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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
	RealWorth         string
	AccumulativeWorth string
	FundCode          string `xorm:"index"`
}

func NewFundInfo(code string, name string, style string) (obj *FundInfo) {
	obj = new(FundInfo)
	obj.Code = code
	obj.Name = name
	obj.Style = style
	return
}

func NewWorth(date string, realWorth string, accWorth string, code string) (obj *Worth) {
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
	if err != nil {
		fmt.Println(err)
		panic(err)
		return
	}
	defer orm.Close()
	err = sync(orm)
	if err != nil {
		panic(err)
		fmt.Println(err)
		return
	}
	baseData, _ := getBaseFundInfoByTT()
	orm.Insert(baseData)
	for _, obj := range baseData {
		worthList, _ := getFundHistory(obj.Code)
		orm.Insert(worthList)
	}
	fmt.Println("update......ok")

}
func InitFundHistory(baseData []*FundInfo, orm *xorm.Engine) error {

	for _, obj := range baseData {
		worthData, err := getFundHistory(obj.Code)
		if err != nil {
			return err
		}
		orm.Insert(worthData)
	}
	return nil
}

func getFundHistory(code string) ([]*Worth, error) {

	url := "http://fund.eastmoney.com/f10/F10DataApi.aspx?"
	datatype := "type=lsjz&"
	fundCode := "code=" + code + "&"
	page := "page=1&"
	per := "per=50Z&"
	// sdate :="sdate=&"
	// edate := "edate=&"
	// rt := "rt=0.04558350201064887"

	reqUrl := url + datatype + fundCode + page + per
	fmt.Println(reqUrl)
	req, err := http.Get(reqUrl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	reqBody := string(body)
	html := strings.Split(reqBody, "\"")[1]
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	worthData := make([]*Worth, 0)
	doc.Find("table tbody tr").Each(func(i int, sl *goquery.Selection) {

		date := ""
		realworth := ""
		accWorth := ""
		sl.Find("td").Each(func(k int, sel *goquery.Selection) {
			if k == 0 {
				date = sel.Text()
				return
			}
			if k == 1 {
				realworth = sel.Text()
				return
			}
			if k == 2 {
				accWorth = sel.Text()
				return
			}
		})
		obj := NewWorth(date, realworth, accWorth, code)
		worthData = append(worthData, obj)
	})
	return worthData, nil
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
				return
			}
			if key == 4 {
				name = sel.Find("nobr a").First().Text()
				return
			}
		})
		obj := NewFundInfo(code, cd.ConvString(name), "")
		baseData = append(baseData, obj)
	})

	return baseData, nil
}
