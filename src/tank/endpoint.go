// Endpoint
// 端点
//
// @author xiangqian
// @date 21:09 2022/12/01
package tank

import (
	"log"
	"net"
)

// Datagram Packet Type
type (
	Ty int8
)

const (
	TyReg    Ty = iota // 注册
	TyDiscov           // 发现
	TyHb               // 心跳
	TyData             // 数据
)

// Datagram Packet
type DgPkt struct {
	Ty   Ty     `json:"ty"`   // 类型
	Data []byte `json:"data"` // 数据
}

type Endpoint struct {
	conn      *net.UDPConn   // 当前端点连接
	localAddr *net.UDPAddr   // 本地地址
	addrs     []*net.UDPAddr // 本地&远程端点udp地址集切片
}

func (endpoint *Endpoint) Listen() {

	// UDP Listen
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0})
	if err != nil {
		panic(err)
	}

	// 关闭连接
	defer conn.Close()

	endpoint.conn = conn
	port := conn.LocalAddr().(*net.UDPAddr).Port
	localAddr := &net.UDPAddr{IP: LocalIp(), Port: port}
	endpoint.localAddr = localAddr
	endpoint.addrs = append(endpoint.addrs, localAddr)
	log.Printf("localAddr: %v\n", localAddr.String())

	count := 0
	var buf [2048]byte
	for {
		// 读取数据
		n, addr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			if count > 8 {
				panic(err)
			}

			log.Printf("read err, %v\n", err)
			count++
			continue
		}

		pDgPkt := &DgPkt{}
		err = Deserialize(buf[:n], pDgPkt)
		if err != nil {
			continue
		}

		for i, length := 0, len(endpoint.addrs); i < length; i++ {

		}

		if pDgPkt.Ty == TyReg {

		}

		log.Printf("addr: %v, data: %v\n", addr, string(buf[:n]))
	}
}

func (endpoint *Endpoint) Write(data []byte) {
	for i, length := 0, len(endpoint.addrs); i < length; i++ {
		addr := endpoint.addrs[i]
		if addr == endpoint.localAddr {
			continue
		}

		// 写入数据
		_, err := endpoint.conn.WriteToUDP(data, addr)
		if err != nil {
			log.Printf("write err, %v\n", err)
		}
	}
}

func LocalIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	return conn.LocalAddr().(*net.UDPAddr).IP
}
