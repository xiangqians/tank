// Endpoint
// 端点
//
// @author xiangqian
// @date 21:09 2022/12/01
package tank

import (
	"log"
	"net"
	"strings"
)

// Datagram Packet Type
type (
	DgPktTy int8
)

const (
	DgPktTyReg    DgPktTy = iota + 1 // 注册
	DgPktTyDiscov                    // 发现
	DgPktTyHb                        // 心跳
	DgPktTyData                      // 数据
)

// Datagram Packet
type DgPkt struct {
	DgPktTy DgPktTy `json:"dgPktTy"` // 类型
	Data    []byte  `json:"data"`    // 数据
}

type Endpoint struct {
	Conn      *net.UDPConn   // 当前端点连接
	LocalAddr *net.UDPAddr   // 本地地址
	Addrs     []*net.UDPAddr // 本地&远程端点udp地址集切片
}

func (pEndpoint *Endpoint) Listen() {
	// UDP Listen
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0})
	if err != nil {
		panic(err)
	}

	// 关闭连接
	defer conn.Close()

	pEndpoint.Conn = conn
	port := conn.LocalAddr().(*net.UDPAddr).Port
	localAddr := &net.UDPAddr{IP: LocalIp(), Port: port}
	pEndpoint.LocalAddr = localAddr
	pEndpoint.Addrs = append(pEndpoint.Addrs, localAddr)
	log.Printf("localAddr: %v\n", localAddr.String())
	pApp.pReg.LocalAddr = localAddr.String()

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
		if err != nil ||
			!(pDgPkt.DgPktTy == DgPktTyReg || pDgPkt.DgPktTy == DgPktTyDiscov || pDgPkt.DgPktTy == DgPktTyHb || pDgPkt.DgPktTy == DgPktTyData) {
			continue
		}

		switch pDgPkt.DgPktTy {
		case DgPktTyReg:
			pEndpoint.Addrs = append(pEndpoint.Addrs, addr)
			pDgPkt := &DgPkt{}
			pDgPkt.DgPktTy = DgPktTyDiscov
			addrs := pEndpoint.Addrs
			var addrsStr string
			for _, a := range addrs {
				addrsStr += a.String() + ","
			}
			pDgPkt.Data = []byte(addrsStr)
			log.Printf("addrsStr: %v\n", addrsStr)
			pEndpoint.SendDgPktToAddrs(pDgPkt)

		case DgPktTyDiscov:
			addrsStr := string(pDgPkt.Data)
			log.Printf("Discov: %v\n", addrsStr)
			addrStrArr := strings.Split(addrsStr, ",")
			for i, as := range addrStrArr {
				as = strings.TrimSpace(as)
				if as == "" {
					continue
				}
				log.Printf("%v, %v\n", i, as)

				flag := true
				for _, _as := range pEndpoint.Addrs {
					if _as.String() == as {
						flag = false
						break
					}
				}
				if flag {
					asAddr, _ := net.ResolveUDPAddr("udp", as)
					pEndpoint.Addrs = append(pEndpoint.Addrs, asAddr)
					pDgPkt := &DgPkt{}
					pDgPkt.DgPktTy = DgPktTyHb
					log.Printf("hb: %v\n", asAddr.String())
					pEndpoint.SendDgPkt(pDgPkt, asAddr)
				}
			}

		case DgPktTyHb:
		case DgPktTyData:
			pAbsGraphics := &AbsGraphics{}
			err := Deserialize(pDgPkt.Data, pAbsGraphics)
			if err == nil {
				log.Printf("addr: %v, graphics: %v\n", addr, pAbsGraphics)
				switch pAbsGraphics.GraphicsTy {
				case GraphicsTyTank:
					pTank := &Tank{AbsGraphics: pAbsGraphics}
					pTank.Init(pTank)
					pApp.pGame.AddGraphics(pTank)
				case GraphicsTyBullet:
				}
			}
		}

		for i, length := 0, len(pEndpoint.Addrs); i < length; i++ {
		}
	}
}

func (pEndpoint *Endpoint) SendGraphics(graphics Graphics) {
	buf, err := Serialize(graphics)
	if err != nil {
		return
	}

	pDgPkt := &DgPkt{}
	pDgPkt.DgPktTy = DgPktTyData
	pDgPkt.Data = buf
	pEndpoint.SendDgPktToAddrs(pDgPkt)
}

func (pEndpoint *Endpoint) SendDgPktToAddrs(pDgPkt *DgPkt) {
	buf, err := Serialize(pDgPkt)
	if err != nil {
		return
	}

	for i, length := 0, len(pEndpoint.Addrs); i < length; i++ {
		addr := pEndpoint.Addrs[i]
		if UDPAddrEqual(addr, pEndpoint.LocalAddr) {
			continue
		}

		_, err := pEndpoint.Write(buf, addr)
		if err != nil {
			log.Printf("write err, %v\n", err)
		}
	}
}

func (pEndpoint *Endpoint) SendDgPkt(pDgPkt *DgPkt, addr *net.UDPAddr) {
	buf, err := Serialize(pDgPkt)
	if err == nil {
		pEndpoint.Write(buf, addr)
	}
}

// 写入数据
func (pEndpoint *Endpoint) Write(data []byte, addr *net.UDPAddr) (int, error) {
	return pEndpoint.Conn.WriteToUDP(data, addr)
}

func UDPAddrEqual(v1, v2 *net.UDPAddr) bool {
	if v1.Port != v2.Port {
		return false
	}

	ip1, ip2 := v1.IP, v1.IP
	if len(ip1) != len(ip2) {
		return false
	}

	for i, length := 0, len(ip1); i < length; i++ {
		if ip1[i] != ip2[i] {
			return false
		}
	}

	return true
}

func LocalIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	return conn.LocalAddr().(*net.UDPAddr).IP
}
