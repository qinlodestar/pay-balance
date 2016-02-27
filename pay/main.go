package main

import (
	log "code.google.com/p/log4go"
	"flag"
	"fmt"
)

func main() {
	fmt.Printf("pay start\n")
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err)
	}
	log.LoadConfiguration(Conf.Log)
	fmt.Printf("%#v\n", Conf)
	initBalance()
	initHttp()
}
