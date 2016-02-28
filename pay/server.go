package main

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"fmt"
	define "github.com/qinlodestar/pay-balance/define"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"strconv"
)

var (
	Db *leveldb.DB
)

type Response struct {
	ErrorCode int32
}

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
		key       string
		msgKey    string
		err       error
		errorCode int32
		status    []byte
	)

	params := r.URL.Query()
	sUserId := params.Get("userId")
	sOrderId := params.Get("orderId")
	sMoney := params.Get("money")
	key = "pay_" + sUserId + "_" + sOrderId
	log.Debug("userId=%s\torderId=%s\tmoney=%s", sUserId, sOrderId, sMoney)

	//如果已经处理成功过的请求就不再处理，直接返回成功
	msgKey = "msg_" + key
	status, err = Db.Get([]byte(msgKey), nil)
	sStatus := fmt.Sprintf("%s", status)
	if sStatus == "1" {
		networkWrite(w, define.SUCCESS)
		log.Debug("success\tkey=%s", key)
		return
	}

	err = Db.Put([]byte(key), []byte(sMoney), nil)
	if err != nil {
		networkWrite(w, define.ERROR_LEVELDB)
		return
	}
	iUserId, _ := strconv.ParseInt(sUserId, 10, 64)
	fMoney, _ := strconv.ParseFloat(sMoney, 64)
	errorCode, err = pushMsg(iUserId, sOrderId, fMoney)
	err = ResponseWrite(w, key, errorCode)
	if err != nil {
		log.Error("trans2balance\tkey=%s\terr=%v", key, err)
		return
	}
}

//只处理网络请求
func networkWrite(w http.ResponseWriter, errorCode int32) {
	var bRes []byte
	var res Response
	res.ErrorCode = errorCode
	bRes, _ = json.Marshal(&res)
	w.Write(bRes)
}

//返回值
func ResponseWrite(w http.ResponseWriter, key string, errorCode int32) (err error) {
	var bRes []byte
	var res Response
	res.ErrorCode = errorCode
	bRes, _ = json.Marshal(&res)
	_, err = w.Write(bRes)
	if err != nil {
		log.Debug("ResponseWrite\tkey=%s\terr=%v", key, err)
	}
	err = Db.Delete([]byte(key), nil)
	if err != nil {
		return
	}
	var status []byte = []byte("1")
	if errorCode != define.SUCCESS {
		status = []byte("0")
	}
	msgKey := "msg_" + key
	err = Db.Put([]byte(msgKey), status, nil)
	if err != nil {
		return
	}
	return
}
