package ping

import (
	"bytes"
	"encoding/binary"
	"net"
	"os/exec"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

type ReqInfo struct {
	DesIp string //目的Ip
	Count int    //请求次数
	Size int    //发送缓冲区大小,单位：字节
	Timeout int64   //等待每次回复的超时时间
	Nerverstop bool  //Ping 指定的主机，直到停止
}
type ResInfo struct {
    SendN int   //发送数据包
    RecvN int   //接收数据包
    LostN int    //丢失数据包
    MinTime int   //往返最短时间
    MaxTime int    //往返最长时间
	TotalTime int   //总共时间
    AvgTime int    //往返平均时间
}

var hosts = []string{
	"www.horizon.ai/",
	"www.baidu.com",
	"www.sina.com",
	"www.taobao.com",
}

/*
当前网络是否连通
true 连通
false 断网
 */
func NetworkConnectivity() bool{
	ans := false
	for _,hostIp :=range hosts {
		//cmd := exec.Command("ping", hostIp, "-c", "1", "-W", "5")
		cmd := exec.Command("ping", hostIp,"-w","1")
		//fmt.Println("NetWorkStatus Start:", time.Now().Unix())
		err := cmd.Run()
		//fmt.Println("NetWorkStatus End  :", time.Now().Unix())
        if err == nil {
        	ans = true
			break
		}
	}
	return ans
}

/*
   返回当前主机到目的Ip的ping信息
 */

func (reqInfo *ReqInfo)NetworkStatus() (*ResInfo,bool){
	resInfo := &ResInfo{
		0,
		0,
		0,
		10000, //这个值一定超时
		0,
		0,
		0,
	}
	if reqInfo == nil {
		return resInfo,false
	}
	conn, err := net.DialTimeout("ip:icmp", reqInfo.DesIp, time.Duration(reqInfo.Timeout) * time.Millisecond)
	if err != nil {
		return resInfo,false
	}
	defer conn.Close()
	icmp := ICMP{
		8,
		0,
		0,
		1,
		1,
	}
	//fmt.Printf("\n正在 ping %s 具有 %d 字节的数据:\n", reqInfo.DesIp, reqInfo.Size)

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp) // 以大端模式写入
	data := make([]byte,  reqInfo.Size)                    //
	buffer.Write(data)
	data = buffer.Bytes()
	count := reqInfo.Count
	for count > 0 || reqInfo.Nerverstop {
		count--
		icmp.SequenceNum = uint16(1)
		// 检验和设为0
		data[2] = byte(0)
		data[3] = byte(0)

		data[6] = byte(icmp.SequenceNum >> 8)
		data[7] = byte(icmp.SequenceNum)
		icmp.Checksum = CheckSum(data)
		data[2] = byte(icmp.Checksum >> 8)
		data[3] = byte(icmp.Checksum)

		// 开始时间
		t1 := time.Now()
		conn.SetDeadline(t1.Add(time.Duration(time.Duration(reqInfo.Timeout) * time.Millisecond)))
		_, err := conn.Write(data)
		if err != nil {
			return resInfo,false
			//log.Fatal(err)
		}
		buf := make([]byte, 65535)
		_, err = conn.Read(buf)
		resInfo.SendN++
		if err != nil {
			//fmt.Println("请求超时。")
			resInfo.LostN++
			continue
		}
		et := int(int64(time.Since(t1) / 1000000))
		if resInfo.MinTime > et {
			resInfo.MinTime = et
		}
		if resInfo.MaxTime <et {
			resInfo.MaxTime = et
		}
		resInfo.TotalTime += et
		//fmt.Printf("来自 %s 的回复: 字节=%d 时间=%dms TTL=%d\n", reqInfo.DesIp, len(buf[28:n]), et, buf[8])
		//resInfo.SendN++
		time.Sleep(1 * time.Second)
	}
	resInfo.RecvN = resInfo.SendN-resInfo.LostN
	if resInfo.RecvN == 0 {
		resInfo.AvgTime = 0
	}else {
		resInfo.AvgTime = resInfo.TotalTime/resInfo.RecvN
	}
	/*
	fmt.Printf("\n%s 的 Ping 统计信息:\n", reqInfo.DesIp)
	fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.2f%% 丢失)，\n", resInfo.SendN, resInfo.RecvN, resInfo.LostN, float64(resInfo.LostN * 100) / float64(resInfo.SendN))
	if resInfo.MaxTime != 0 && resInfo.MinTime != int(math.MaxInt32) {
		fmt.Printf("往返行程的估计时间(以毫秒为单位):\n")
		//fmt.Println("totalTIme)
		fmt.Printf("    最短 = %dms，最长 = %dms，平均 = %dms\n", resInfo.MinTime, resInfo.MaxTime, resInfo.TotalTime / resInfo.RecvN)
	}
	*/
	return resInfo,true
}

func CheckSum(data []byte) uint16 {
	var sum uint32
	var length = len(data)
	var index int
	for length > 1 { // 溢出部分直接去除
		sum += uint32(data[index]) << 8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length == 1 {
		sum += uint32(data[index])
	}
	// CheckSum的值是16位，计算是将高16位加低16位，得到的结果进行重复以该方式进行计算，直到高16位为0
	/*
		sum的最大情况是：ffffffff
		第一次高16位+低16位：ffff + ffff = 1fffe
		第二次高16位+低16位：0001 + fffe = ffff
		即推出一个结论，只要第一次高16位+低16位的结果，再进行之前的计算结果用到高16位+低16位，即可处理溢出情况
	 */
	sum = uint32(sum >> 16) + uint32(sum)
	sum = uint32(sum >> 16) + uint32(sum)
	return uint16(^sum)
}



