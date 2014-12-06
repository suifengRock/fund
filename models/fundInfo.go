package models

type FundInfo struct {
	Id    int64
	Code  string `xorm:"index unique"`
	Name  string `xorm:"not null"`
	Style string
}

func NewFundInfo(code string, name string, style string) (obj *FundInfo) {
	obj = new(FundInfo)
	obj.Code = code
	obj.Name = name
	obj.Style = style
	return
}
