package dns

import (
	"bytes"
	"os/exec"
	"regexp"
)

/*
   如果可以连上目的IP
   data[0]返回本机配置dns
   data[1]返回目的IP DNS
 */
func DNSGet(DstIp string) (string,string,bool){
	cmd := exec.Command("nslookup", DstIp)
	//cmd.Stdin = strings.NewReader("some input") //输入
	var out bytes.Buffer
	cmd.Stdout = &out //输出
	err := cmd.Run()
	if err != nil {
		//fmt.Println("获取dns失败")
		return "","",false
	}
	reg := regexp.MustCompile("(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}")
	datas := reg.FindAllString(out.String(),2)
	lendatas := len(datas)
    if lendatas < 2{
		return "","",false
	}
	return datas[0],datas[1],true
}
