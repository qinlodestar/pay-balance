package main

import (
	log "code.google.com/p/log4go"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/qinlodestar/pay-balance/libs/mysql"
	//"fmt"
)

func main() {
	//fmt.Printf("pay start\n")
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err)
	}
	log.LoadConfiguration(Conf.Log)
	//fmt.Printf("%#v\n", Conf)
	initBalance()

	var err error
	mysql.MyDb, err = sql.Open("mysql", Conf.MysqlAddr)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer mysql.MyDb.Close()

	initHttp()
}
