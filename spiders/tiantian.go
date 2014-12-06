package spiders

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-xorm/xorm"
	"github.com/qiniu/iconv"
	"github.com/suifengRock/fund/models"
	"io/ioutil"
	"net/http"
	"strings"
)

func InitFundHistory(baseData []*models.FundInfo, orm *xorm.Engine) error {

	for _, obj := range baseData {
		worthData, err := getFundHistory(obj.Code)
		if err != nil {
			return err
		}
		orm.Insert(worthData)
	}
	return nil
}

func GetFundHistory(code string) ([]*Worth, error) {

	url := "http://fund.eastmoney.com/f10/F10DataApi.aspx?"
	datatype := "type=lsjz&"
	fundCode := "code=" + code + "&"
	page := "page=1&"
	per := "per=50&"
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
	worthData := make([]*models.Worth, 0)
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

func GetBaseFundInfoByTT() ([]*FundInfo, error) {

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
