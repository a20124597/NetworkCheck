package main

import (
	"NetworkCheck/dns"
	"NetworkCheck/ping"
	"fmt"
)

func main() {
	//test newworkd connectivity
    fmt.Println(
    	` _    _            _                  _   _      _     _______          _     
| |  | |          (_)                | \ | |    | |   |__   __|        | |    
| |__| | ___  _ __ _ _______  _ __   |  \| | ___| |_     | | ___   ___ | |___ 
|  __  |/ _ \|  __| |_  / _ \|  _ \  | .   |/ _ \ __|    | |/ _ \ / _ \| / __|
| |  | | (_) | |  | |/ / (_) | | | | | |\  |  __/ |_     | | (_) | (_) | \__ \
|_|  |_|\___/|_|  |_/___\___/|_| |_| |_| \_|\___|\__|    |_|\___/ \___/|_|___/`)
    fmt.Println("Begin To Check Network Connectivity...")
	networkConnectivity := ping.NetworkConnectivity()
	if networkConnectivity {
		fmt.Println("NetWork Stats OK!")
	}else {
		fmt.Println("NetWork Stats Wrong!")
	}

	fmt.Println("Begin To Check DNS Config...")
	reqInfo := &ping.ReqInfo{
		"www.baidu.com",
		4,
		32,
		1000,
		false,
	}
	localDns,_,status := dns.DNSGet(reqInfo.DesIp)
	if status {
		fmt.Println("Local DNS Is ",localDns)
	}else {
		fmt.Println("DNS Config Is Wrong!")
	}
	fmt.Println("Begin To Connect www.horizon.ai...")
    reqInfo.DesIp = "www.horizon.ai"
    resInfo,status := reqInfo.NetworkStatus()
    if status {
		fmt.Println("往返行程的估计时间(以毫秒为单位:)")
		fmt.Println("最短 = ",resInfo.MinTime,"ms","最长 = ",resInfo.MaxTime,"ms","平均 =",resInfo.AvgTime,"ms")
	}else {
		fmt.Println("Connect Server Error!")
	}

}
