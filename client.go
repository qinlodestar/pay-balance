package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	//"time"
)

var wg sync.WaitGroup

func get(i int) {
	var (
		res  *http.Response
		err  error
		body []byte
	)
	//url := fmt.Sprintf("http://127.0.0.1:8070/sub?ver=1&op=7&seq=1&cb=callback&t=%d", i)
	url := fmt.Sprintf("http://123.57.72.193:17810/trans2balance?userId=%d&orderId=100&money=99.9", i)
	//fmt.Printf("%s\n", url)
	if res, err = http.Get(url); err != nil {
		fmt.Printf("%d\n%v\n", i, err)
	}
	body, err = ioutil.ReadAll(res.Body)
	fmt.Printf("%s", body)
	defer wg.Done()
}

func main() {
	for i := 1; i < 2; i++ {
		//if i%50 == 0 {
		wg.Add(1)
		//	time.Sleep(time.Second * 1)
		//}
		go get(i)
	}
	wg.Wait()
	//time.Sleep(time.Second * 10000)
}
