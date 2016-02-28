package main

import (
	log "code.google.com/p/log4go"
	"flag"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

var (
	WG sync.WaitGroup
	Db *leveldb.DB
)

func main() {
	var err error
	fmt.Printf("balance start")
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err)
	}
	log.LoadConfiguration(Conf.Log)
	fmt.Printf("%#v\n", Conf)

	Db, err = leveldb.OpenFile("/data/leveldb/balance", nil)
	if err != nil {
		panic("db open failed")
	}
	//initHttp()
	initRpc()
	WG.Add(1)
	WG.Wait()
}
