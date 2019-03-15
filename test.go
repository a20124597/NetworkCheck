package main

import (
	"net"
	"fmt"
	"time"
)

func main() {
	hostname := "www.sina.com"
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		fmt.Printf("lookup host error: %v\n", err)
	} else {
		fmt.Printf("addrs: %v", addrs)
	}
	t :=time.Duration(1) * time.Millisecond
	println(t)
	t1 :=time.Duration(1000)* 1000*1000
	println(t1)
}