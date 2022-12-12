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
	pConn      *net.UDPConn   // 当前端点连接
	pLocalAddr *net.UDPAddr   // 本地地址
	pAddrs     []*net.UDPAddr // 本地&远程端点udp地址集切片
}

func (pEndpoint *Endpoint) Listen() {
	// UDP Listen
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: DefaultEndpointPort})
	if err != nil {
		panic(err)
	}

	// 关闭连接
	defer conn.Close()

	pEndpoint.pConn = conn
	port := conn.LocalAddr().(*net.UDPAddr).Port
	localAddr := &net.UDPAddr{IP: LocalIp(), Port: port}
	pEndpoint.pLocalAddr = localAddr
	pEndpoint.pAddrs = append(pEndpoint.pAddrs, localAddr)
	log.Printf("localAddr: %v\n", localAddr.String())
	pApp.pReg.SetLocalAddr(localAddr.String())

	count := 0
	var buf [2048]byte
	for {
		// 读取数据
		n, pAddr, err := conn.ReadFromUDP(buf[:])
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
			pEndpoint.receiveRegDgPkt(pDgPkt, pAddr)

		case DgPktTyDiscov:
			pEndpoint.receiveDiscovDgPkt(pDgPkt, pAddr)

		case DgPktTyHb:
			pEndpoint.receiveHbDgPkt(pDgPkt, pAddr)

		case DgPktTyData:
			pEndpoint.receiveDataDgPkt(pDgPkt, pAddr)
		}
	}
}

// 接收到注册数据包
func (pEndpoint *Endpoint) receiveRegDgPkt(pDgPkt *DgPkt, pAddr *net.UDPAddr) {
	log.Printf("receiveRegDgPkt: %v\n", pAddr.String())

	// 添加到 本地&远程端点udp地址集切片
	pEndpoint.pAddrs = append(pEndpoint.pAddrs, pAddr)

	// 反序列注册者坦克信息
	graphics := DeserializeBytesToGraphics(pDgPkt.Data)
	if graphics != nil {
		// 将注册者坦克信息添加到图形集
		pApp.pGame.AddGraphics(graphics)

		// 让所有端都发现新上线坦克
		pEndpoint.SendGraphicsToAddrs(graphics)
	}
	//log.Printf("reg AbsGraphics: %v\n", *pAbsGraphics)

	// SendDiscovDgPktToAddrs
	pDgPkt.DgPktTy = DgPktTyDiscov
	var addrsStr string
	for _, a := range pEndpoint.pAddrs {
		addrsStr += a.String() + ","
	}
	pDgPkt.Data = []byte(addrsStr)
	log.Printf("SendDiscovDgPktToAddrs: %v\n", addrsStr)
	pEndpoint.SendDgPktToAddrs(pDgPkt)

	// 阻塞获取 chanel 中的 map
	graphicsMap := <-pApp.pGame.GraphicsMapChan

	// 再将 map 添加到 channel
	defer func() { pApp.pGame.GraphicsMapChan <- graphicsMap }()

	// 发送当前图形集给注册者
	for _, graphics := range graphicsMap {
		pEndpoint.SendGraphics(graphics, pAddr)
	}
}

// 接收到发现数据包
func (pEndpoint *Endpoint) receiveDiscovDgPkt(pDgPkt *DgPkt, pAddr *net.UDPAddr) {
	addrsStr := string(pDgPkt.Data)
	log.Printf("receiveDiscovDgPkt: %v, %v\n", pAddr.String(), addrsStr)
	addrsStrArr := strings.Split(addrsStr, ",")
	for i, as := range addrsStrArr {
		as = strings.TrimSpace(as)
		if as == "" {
			continue
		}
		log.Printf("%v, %v\n", i, as)

		flag := true
		for _, _as := range pEndpoint.pAddrs {
			if _as.String() == as {
				flag = false
				break
			}
		}
		if flag {
			asAddr, _ := net.ResolveUDPAddr("udp", as)
			pEndpoint.pAddrs = append(pEndpoint.pAddrs, asAddr)
			pDgPkt.DgPktTy = DgPktTyHb
			pDgPkt.Data = nil
			//log.Printf("hb: %v\n", asAddr.String())
			pEndpoint.SendDgPkt(pDgPkt, asAddr)
		}
	}
}

// 接收到心跳数据包
func (pEndpoint *Endpoint) receiveHbDgPkt(pDgPkt *DgPkt, pAddr *net.UDPAddr) {
	log.Printf("receiveHbDgPkt: %v\n", pAddr.String())
}

// 接收到数据包
func (pEndpoint *Endpoint) receiveDataDgPkt(pDgPkt *DgPkt, pAddr *net.UDPAddr) {
	graphics := DeserializeBytesToGraphics(pDgPkt.Data)
	if graphics != nil {
		pApp.pGame.AddGraphics(graphics)
	}
}

func (pEndpoint *Endpoint) SendGraphicsToAddrs(graphics Graphics) {
	if graphics == nil {
		return
	}

	buf, err := Serialize(graphics)
	if err != nil {
		return
	}

	pDgPkt := &DgPkt{}
	pDgPkt.DgPktTy = DgPktTyData
	pDgPkt.Data = buf
	pEndpoint.SendDgPktToAddrs(pDgPkt)
}

func (pEndpoint *Endpoint) SendGraphics(graphics Graphics, pAddr *net.UDPAddr) {
	if graphics == nil {
		return
	}

	buf, err := Serialize(graphics)
	if err != nil {
		return
	}

	pDgPkt := &DgPkt{}
	pDgPkt.DgPktTy = DgPktTyData
	pDgPkt.Data = buf
	pEndpoint.SendDgPkt(pDgPkt, pAddr)
}

func (pEndpoint *Endpoint) SendDgPktToAddrs(pDgPkt *DgPkt) {
	buf, err := Serialize(pDgPkt)
	if err != nil {
		return
	}

	// debug
	//str := ""
	//var count = 0

	for i, length := 0, len(pEndpoint.pAddrs); i < length; i++ {
		pAddr := pEndpoint.pAddrs[i]
		if UDPAddrEqual(pAddr, pEndpoint.pLocalAddr) {
			continue
		}

		_, err := pEndpoint.Write(buf, pAddr)
		if err != nil {
			log.Printf("write err, %v\n", err)
		}

		// debug
		//if str != "" {
		//	str += ","
		//}
		//str += pAddr.String()
		//count++
	}

	// debug
	//if count > 0 {
	//	log.Printf("Send to %v endpoints: %v\n", count, str)
	//}
}

func (pEndpoint *Endpoint) SendDgPkt(pDgPkt *DgPkt, pAddr *net.UDPAddr) bool {
	if UDPAddrEqual(pAddr, pEndpoint.pLocalAddr) {
		return false
	}

	buf, err := Serialize(pDgPkt)
	if err == nil {
		_, err = pEndpoint.Write(buf, pAddr)
	}

	return err == nil
}

// 写入数据
func (pEndpoint *Endpoint) Write(data []byte, pAddr *net.UDPAddr) (int, error) {
	return pEndpoint.pConn.WriteToUDP(data, pAddr)
}

func UDPAddrEqual(v1, v2 *net.UDPAddr) bool {
	//if v1.Port != v2.Port {
	//	return false
	//}
	//
	//ip1, ip2 := v1.IP, v2.IP
	//if len(ip1) != len(ip2) {
	//	return false
	//}
	//
	//for i, length := 0, len(ip1); i < length; i++ {
	//	if ip1[i] != ip2[i] {
	//		return false
	//	}
	//}
	//
	//return true

	return v1.String() == v2.String()
}

func LocalIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	return conn.LocalAddr().(*net.UDPAddr).IP
}
