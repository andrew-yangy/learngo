package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (fe *frontendServer) getOrders() string {
	resp, err := http.Get("http://order.default.svc.cluster.local")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}
