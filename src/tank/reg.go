// Reg
// @author xiangqian
// @date 13:12 2022/12/03
package tank

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type InputFlag int8

const (
	InputFlagIp = iota + 1
	InputFlagPort
)

type Reg struct {
	LocalAddr    string
	Ip           string
	Port         string
	ArcadeFont   font.Face
	InputFlag    InputFlag
	CursorStatus bool
}

func (pReg *Reg) Init() {
	pReg.ArcadeFont = CreateFontFace(16, 72)
	pReg.InputFlag = InputFlagIp

	go func() {
		for pApp.appStep == AppStepReg {
			pReg.CursorStatus = !pReg.CursorStatus
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func (pReg *Reg) SetLocalAddr(localAddr string) {
	pReg.LocalAddr = localAddr
	pReg.Ip = strings.Split(pReg.LocalAddr, ":")[0]
}

func (pReg *Reg) Append(a string) {
	// ip
	if pReg.InputFlag == InputFlagIp {
		ip := pReg.Ip
		length := len(ip)
		if length >= len("255.255.255.255") {
			return
		}

		if a == "." && (length == 0 || strings.Count(ip, ".") >= 3) {
			return
		}

		if a == "." && length > 0 && string(ip[length-1]) == a {
			return
		}

		lastIndex := strings.LastIndex(ip, ".")

		if strings.Count(ip, ".") >= 3 && length-lastIndex > 3 {
			return
		}

		if lastIndex > 0 && length-lastIndex <= 3 {
			i, _ := strconv.ParseInt(ip[lastIndex+1:]+a, 10, 64)
			log.Printf("i = %v, %v\n", ip[lastIndex+1:]+a, i)
			if i > 255 {
				return
			}
		}

		if (lastIndex == -1 && length == 3) || (lastIndex > 0 && length-lastIndex > 3) {
			pReg.Ip += "."
		}

		pReg.Ip += a
		return
	}

	// port
	if len(pReg.Port) >= 5 || a == "." {
		return
	}
	pReg.Port += a
}

func (pReg *Reg) Subtract() {
	// ip
	if pReg.InputFlag == InputFlagIp {
		length := len(pReg.Ip)
		if length > 0 {
			pReg.Ip = pReg.Ip[:length-1]
		}
		return
	}

	// port
	length := len(pReg.Port)
	if length > 0 {
		pReg.Port = pReg.Port[:length-1]
	}
}

func (pReg *Reg) SendRegDgPkt() {
	ip := pReg.Ip
	var ipArr []string = nil
	if len(ip) > 0 {
		ipArr = strings.Split(ip, ".")
	}
	if len(pReg.Port) > 0 && len(ip) > 0 && len(ipArr) == 4 {
		pDgPkt := &DgPkt{}
		pDgPkt.DgPktTy = DgPktTyReg

		port, _ := strconv.ParseInt(pReg.Port, 10, 64)
		a, _ := strconv.ParseInt(ipArr[0], 10, 64)
		b, _ := strconv.ParseInt(ipArr[1], 10, 64)
		c, _ := strconv.ParseInt(ipArr[2], 10, 64)
		d, _ := strconv.ParseInt(ipArr[3], 10, 64)
		addr := &net.UDPAddr{
			IP:   net.IPv4(byte(a), byte(b), byte(c), byte(d)),
			Port: int(port),
		}
		log.Printf("input Ip: %v, Port: %v\n", pReg.Ip, pReg.Port)
		log.Printf("reg addr: %v\n", addr.String())
		pApp.pEndpoint.SendDgPkt(pDgPkt, addr)
	}
	pApp.appStep = AppStepGame
}

func (pReg *Reg) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyKP0) || inpututil.IsKeyJustPressed(ebiten.Key0) {
		pReg.Append("0")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP1) || inpututil.IsKeyJustPressed(ebiten.Key1) {
		pReg.Append("1")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP2) || inpututil.IsKeyJustPressed(ebiten.Key2) {
		pReg.Append("2")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP3) || inpututil.IsKeyJustPressed(ebiten.Key3) {
		pReg.Append("3")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP4) || inpututil.IsKeyJustPressed(ebiten.Key4) {
		pReg.Append("4")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP5) || inpututil.IsKeyJustPressed(ebiten.Key5) {
		pReg.Append("5")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP6) || inpututil.IsKeyJustPressed(ebiten.Key6) {
		pReg.Append("6")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP7) || inpututil.IsKeyJustPressed(ebiten.Key7) {
		pReg.Append("7")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP8) || inpututil.IsKeyJustPressed(ebiten.Key8) {
		pReg.Append("8")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP9) || inpututil.IsKeyJustPressed(ebiten.Key9) {
		pReg.Append("9")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		pReg.InputFlag = InputFlagIp ^ InputFlagPort ^ pReg.InputFlag
	} else if inpututil.IsKeyJustPressed(ebiten.KeyPeriod) { // “.” 按键
		pReg.Append(".")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		pReg.Subtract()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		pReg.SendRegDgPkt()
	}

	return nil
}

func (pReg *Reg) Draw(screen *ebiten.Image) {
	x := screenWidth/2 - 230
	y := screenHeight/2 - 200

	text.Draw(screen,
		fmt.Sprintf("localAddr %v", pReg.LocalAddr),
		pReg.ArcadeFont,
		x, y,
		color.White)

	y += 100
	ipY := y
	ipText := fmt.Sprintf("IP:   %v", pReg.Ip)
	text.Draw(screen,
		ipText,
		pReg.ArcadeFont,
		x, y,
		color.White)

	y += 50
	portY := y
	portText := fmt.Sprintf("Port: %v", pReg.Port)
	text.Draw(screen,
		portText,
		pReg.ArcadeFont,
		x, y,
		color.White)

	if pReg.CursorStatus {
		var t string
		if pReg.InputFlag == InputFlagIp {
			t = ipText
			y = ipY
		} else {
			t = portText
			y = portY
		}

		text.Draw(screen,
			fmt.Sprintf("|"),
			pReg.ArcadeFont,
			x+len(t)*16, y,
			color.White)
	}

}
