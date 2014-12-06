package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/suifengRock/fund/models"
	"os"
)

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
	return engine.Sync(&models.FundInfo{}, &models.Worth{})
}
