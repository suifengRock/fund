package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"strings"
)

type FundInfo struct {
	Id    int64
	Code  string
	Name  string
	Style string
}

type Worth struct {
	Id                int64
	Date              string //time.time
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

	conn := "root:123456@/test?charset=utf8"

	return xorm.NewEngine("mysql", conn)
}

func sync(engine *xorm.Engine) error {
	return engine.Sync(&FundInfo{}, &Worth{})
}

func main() {

	// orm, err := mysqlEngine()
	// defer orm.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// 	return
	// }
	// err = sync(orm)
	// if err != nil {
	// 	panic(err)
	// 	fmt.Println(err)
	// 	return
	// }
	fmt.Println("...")

	// tianFund := "http://fund.eastmoney.com/fund.html"
	haoFund := "http://www.howbuy.com/board/"
	// resp, err := http.Get(tianFund)
	// defer res.Body.Close()
	// if err != nil {
	// 	fmt.Println("error: ", err)
	// }
	// b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("11111")
	doc, err := goquery.NewDocument(haoFund)
	if err != nil {
		fmt.Println(err)
		return
	}
	dataList := doc.Find(".result_list")
	dataList.Find("textarea").Each(func(num int, sl *goquery.Selection) {
		if num > 0 {
			return
		}
		sl.RemoveAttr("attr")
		textArea := sl.Text()
		textList := strings.Split(textArea, "<tr>")
		fmt.Println("len", len(textList))
		textSli := strings.Split(textArea, textList[0])
		length := len(textSli)
		for tmp := 0; tmp < length; tmp += 2 {
			val := textSli[tmp]
			valList := strings.Split(val, "</tr>")
			html := valList[0] + "</tr>"
			html = strings.Replace(html, "\"", "'", -1)
			fmt.Println(html)
			dataList.Find("table tbody").AppendHtml("<tr><td width='4%'><input  type='checkbox' onclick='move(this);' value='161605'/></td><td width='5%'>100</td><td width='6%'><a target='_blank' href='/fund/161605/'>161605</a></td><td  class='tdl'><a href='/fund/161605/'>融通蓝筹成长混合</a></td><td  class='tdr'>1.2260</td><td  class='tdr'>2.5590</td><td  class='tdr'>1.2070</td><td  class='tdr'>2.5400</td><td  class='tdr'><span class='cRed'>0.0190</span></td><td  class='tdr'><span class='cRed'>1.57%</span></td><td  class='operate'><a href='https://trade.ehowbuy.com/trade/subs.htm?method=apply&fundCode=161605' class='sg'  title='申购'></a><a href='https://trade.ehowbuy.com/savingplan/index.htm?method=apply&fundCode=161605' class='dt' title='定投'></a><a href='javascript:void(0)' url_add='http://www.howbuy.com/fund/fundtool/addAttentionFund.htm?fundCode=161605' url_del='http://www.howbuy.com/fund/fundtool/removefund.htm?fundCode=161605' class='zx' title='自选' target='_self'></a></td></tr>")
		}

	})

	dataList.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		code := ""
		name := ""
		s.Find("td").Each(func(key int, sel *goquery.Selection) {

			if key == 2 {
				code = sel.Text()
			}
			if key == 3 {
				name = sel.Text()
			}
		})
		fmt.Println(code, " ", name)
	})

	// fmt.Println(dataList.Find("table tbody tr").Text())

	// doc.Find(".reviews-wrap article .review-rhs").Each(func(i int, s *goquery.Selection) {
	// 	band := s.Find("h3").Text()
	// 	title := s.Find("i").Text()
	// 	fmt.Printf("Review %d: %s - %s\n", i, band, title)
	// })

	// html := doc.Find(".tableDiv").Text()
	// fmt.Print(html)
	// doc.Find("tableDiv table tbody tr").Each(func(i int, s *goquery.Selection) {

	// 	if i > 2 {
	// 		fmt.Println(i)
	// 		return
	// 	}
	// 	s.Find("td").Each(func(k int, sel *goquery.Selection) {
	// 		text := sel.Text()
	// 		fmt.Println(text)
	// 	})

	// })

}
