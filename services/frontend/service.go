package main

import (
	"fmt"
	"github.com/ddvkid/learngo/internal/util"
	"io/ioutil"
	"net/http"
)

var (
	namespace = util.GetEnv("K8S_NAMESPACE", "default")
	orderAPI  = fmt.Sprintf("http://order.%s.svc.cluster.local", namespace)
)

func (fe *frontendServer) getOrders() string {
	resp, err := http.Get(orderAPI)
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
