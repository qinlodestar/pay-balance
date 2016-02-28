package main

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"fmt"
	define "github.com/qinlodestar/pay-balance/define"
	mysql "github.com/qinlodestar/pay-balance/libs/mysql"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"strconv"
	"time"
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
	time.Sleep(1 * time.Millisecond)
	//time.Sleep(10000 * time.Millisecond)
	log.Debug("userId=%s\torderId=%s\tmoney=%s", sUserId, sOrderId, sMoney)

	//如果已经处理成功过的请求就不再处理，直接返回成功
	msgKey = "msg_" + key
	status, err = Db.Get([]byte(msgKey), nil)
	sStatus := fmt.Sprintf("%s", status)
	if sStatus == "0" {
		networkWrite(w, define.SUCCESS)
		log.Debug("success\tkey=%s", key)
		return
	}

	iUserId, _ := strconv.ParseInt(sUserId, 10, 64)
	fMoney, _ := strconv.ParseFloat(sMoney, 64)

	//查询资金余额是否满足支付条件
	is := IsEnoughFund(iUserId, fMoney)
	log.Debug("is=%t\n", is)
	if !is {
		networkWrite(w, define.ERROR_FUND_NOT_ENOUGH)
	}

	//写入本地数据库
	err = Db.Put([]byte(key), []byte(sMoney), nil)
	if err != nil {
		networkWrite(w, define.ERROR_LEVELDB)
		return
	}
	errorCode, err = pushMsg(iUserId, sOrderId, fMoney)

	//提交数据库
	is = commitDb(sUserId, sMoney)
	if !is {
		networkWrite(w, define.ERROR_MYSQL)
		return
	}

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
	log.Debug("ResponseWrite\tSUCCESS\tkey=%s", key)
	return
}

func IsEnoughFund(iUserId int64, fMoney float64) (is bool) {
	str := fmt.Sprintf("select * from pay where userId=%d\n", iUserId)
	fmt.Printf("%s", str)
	result := mysql.Query(str)
	//fmt.Printf("%v\n", result)
	length := len(result)
	log.Debug("length=%d", length)
	if length == 0 {
		is = false
		return
	}
	sAllMoney := result[0]["money"]
	log.Debug("sAllMoney=%s", sAllMoney)
	fAllMoney, _ := strconv.ParseFloat(sAllMoney, 64)
	log.Debug("fAllMoney=%f", fAllMoney)
	log.Debug("fMoney=%f", fMoney)

	if fAllMoney < fMoney {
		is = false
		return
	}
	is = true
	return
}

func commitDb(sUserId string, sMoney string) (is bool) {
	is = false
	_, affectRow, _ := mysql.Exec("update pay set money=money-? where userId=?", sUserId, sMoney)
	if affectRow > 0 {
		is = true
	}
	return
}
