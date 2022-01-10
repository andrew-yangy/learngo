package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	namespace = OrString(os.Getenv("K8S_NAMESPACE"), "default")
	orderAPI  = fmt.Sprintf("http://order.%s.svc.cluster.local", namespace)
)

func (fe *frontendServer) getOrders() string {
	fmt.Print(namespace)
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
