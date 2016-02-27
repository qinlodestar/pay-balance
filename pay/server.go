package main

import (
	log "code.google.com/p/log4go"
	//	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"strconv"
)

var (
	Db *leveldb.DB
)

func initHttp() {
	var err error
	log.Debug("starting server...")
	http.HandleFunc("/trans2balance", trans2balance)
	Db, err = leveldb.OpenFile("/data/leveldb/pay", nil)
	if err != nil {
		panic("db open failed")
	}
	defer Db.Close()
	http.ListenAndServe(Conf.HttpBind, nil)
}

func trans2balance(w http.ResponseWriter, r *http.Request) {
	var (
		key string
		err error
	)
	params := r.URL.Query()
	sUserId := params.Get("userId")
	sOrderId := params.Get("orderId")
	sMoney := params.Get("money")
	key = "pay_" + sUserId + "_" + sOrderId
	log.Debug("userId=%s\torderId=%s\tmoney=%s", sUserId, sOrderId, sMoney)
	err = Db.Put([]byte(key), []byte("123"), nil)
	if err != nil {
		_, err = w.Write([]byte("{\"errorCode\":40000}"))
		return
	}
	_, err = w.Write([]byte("{\"errorCode\":0}"))
	if err != nil {
		//删除db数据，
		Db.Delete([]byte(key), nil)
		//如果再删除失败，就需要再次删除，直到成功为止
	}
	//data, err := Db.Get([]byte(key), nil)
	//fmt.Printf("%s\n", data)
	iUserId, _ := strconv.ParseInt(sUserId, 10, 64)
	fMoney, _ := strconv.ParseFloat(sMoney, 64)
	pushMsg(iUserId, sOrderId, fMoney)
}
