package main

import (
	log "code.google.com/p/log4go"
	"flag"
	"fmt"
	"sync"
)

var WG sync.WaitGroup

func main() {
	fmt.Printf("balance start")
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err)
	}
	log.LoadConfiguration(Conf.Log)
	fmt.Printf("%#v\n", Conf)
	//initHttp()
	initRpc()
	WG.Add(1)
	WG.Wait()
}
