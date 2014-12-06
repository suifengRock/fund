package models

type Worth struct {
	Id                int64
	Date              string
	RealWorth         string
	AccumulativeWorth string
	FundCode          string `xorm:"index"`
}

func NewWorth(date string, realWorth string, accWorth string, code string) (obj *Worth) {
	obj = new(Worth)
	obj.Date = date
	obj.RealWorth = realWorth
	obj.AccumulativeWorth = accWorth
	obj.FundCode = code
	return
}
