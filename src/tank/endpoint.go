// Endpoint
// 端点
//
// @author xiangqian
// @date 21:09 2022/12/01
package tank

import (
	"fmt"
	"log"
	"net"
)

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
	endpoint.localAddr = &net.UDPAddr{IP: LocalIp(), Port: port}
	log.Printf("localAddr: %v\n", endpoint.localAddr.String())

	count := 0
	var bf [2048]byte
	for {
		// 读取数据
		n, addr, err := conn.ReadFromUDP(bf[:])
		if err != nil {
			if count > 8 {
				panic(err)
			}

			fmt.Printf("read err, %v\n", err)
			count++
			continue
		}

		fmt.Printf("addr: %v, data: %v\n", addr, string(bf[:n]))
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
			fmt.Printf("write err, %v\n", err)
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
