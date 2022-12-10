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
	InputFlagName = iota + 1
	InputFlagIp
	InputFlagPort
)

type Reg struct {
	LocalAddr    string
	Name         string
	Ip           string
	Port         string
	ArcadeFont   font.Face
	TitleFont    font.Face
	InputFlag    InputFlag
	CursorStatus bool
}

func (pReg *Reg) Init() {
	pReg.ArcadeFont = CreateFontFace(16, 72)
	pReg.TitleFont = CreateFontFace(64, 72)
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

func (pReg *Reg) AppendName(a string) {
	if pReg.InputFlag != InputFlagName {
		return
	}

	if len(pReg.Name) >= 8 {
		return
	}

	//pReg.Name += a
	pReg.Name += strings.ToUpper(a)
}

func (pReg *Reg) AppendIp(a string) {
	if pReg.InputFlag != InputFlagIp {
		return
	}

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
		//log.Printf("i = %v, %v\n", ip[lastIndex+1:]+a, i)
		if i > 255 {
			return
		}
	}

	if a != "." && ((lastIndex == -1 && length == 3) || (lastIndex > 0 && length-lastIndex > 3)) {
		pReg.Ip += "."
	}

	pReg.Ip += a
}

func (pReg *Reg) AppendPort(a string) {
	if pReg.InputFlag != InputFlagPort {
		return
	}

	if len(pReg.Port) >= 5 {
		return
	}

	pReg.Port += a
}

func (pReg *Reg) Subtract() {
	switch pReg.InputFlag {
	case InputFlagName:
		length := len(pReg.Name)
		if length > 0 {
			pReg.Name = pReg.Name[:length-1]
		}

	case InputFlagIp:
		length := len(pReg.Ip)
		if length > 0 {
			pReg.Ip = pReg.Ip[:length-1]
		}

	case InputFlagPort:
		length := len(pReg.Port)
		if length > 0 {
			pReg.Port = pReg.Port[:length-1]
		}
	}
}

func (pReg *Reg) SendRegDgPkt() {
	if pReg.Name == "" {
		pReg.InputFlag = InputFlagName
		return
	}

	pApp.pGame.pTank.Name = pReg.Name

	ip := pReg.Ip
	var ipArr []string = nil
	if len(ip) > 0 {
		ipArr = strings.Split(ip, ".")
	}
	if len(pReg.Port) > 0 && len(ip) > 0 && len(ipArr) == 4 {
		pDgPkt := &DgPkt{}
		pDgPkt.DgPktTy = DgPktTyReg
		buf, err := Serialize(pApp.pGame.pTank)
		if err == nil {
			pDgPkt.Data = buf
		}

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
		a := "0"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP1) || inpututil.IsKeyJustPressed(ebiten.Key1) {
		a := "1"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP2) || inpututil.IsKeyJustPressed(ebiten.Key2) {
		a := "2"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP3) || inpututil.IsKeyJustPressed(ebiten.Key3) {
		a := "3"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP4) || inpututil.IsKeyJustPressed(ebiten.Key4) {
		a := "4"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP5) || inpututil.IsKeyJustPressed(ebiten.Key5) {
		a := "5"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP6) || inpututil.IsKeyJustPressed(ebiten.Key6) {
		a := "6"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP7) || inpututil.IsKeyJustPressed(ebiten.Key7) {
		a := "7"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP8) || inpututil.IsKeyJustPressed(ebiten.Key8) {
		a := "8"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP9) || inpututil.IsKeyJustPressed(ebiten.Key9) {
		a := "9"
		pReg.AppendName(a)
		pReg.AppendIp(a)
		pReg.AppendPort(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		a := "a"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		a := "b"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		a := "c"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		a := "d"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		a := "e"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		a := "f"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		a := "g"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		a := "h"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		a := "i"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		a := "j"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		a := "k"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		a := "l"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		a := "m"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		a := "n"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		a := "o"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		a := "p"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		a := "q"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		a := "r"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		a := "s"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		a := "t"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		a := "u"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		a := "v"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		a := "w"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		a := "x"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyY) {
		a := "y"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		a := "z"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyMinus) { // - 减号
		a := "-"
		pReg.AppendName(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyPeriod) { // “.” 键
		a := "."
		pReg.AppendName(a)
		pReg.AppendIp(a)

	} else if inpututil.IsKeyJustPressed(ebiten.KeyTab) { // tab 键
		pReg.InputFlag++
		if pReg.InputFlag > InputFlagPort {
			pReg.InputFlag = InputFlagName
		}

	} else if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) { // Backspace 键
		pReg.Subtract()

	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) { // Enter 键
		pReg.SendRegDgPkt()
	}

	return nil
}

func (pReg *Reg) Draw(screen *ebiten.Image) {
	x := screenWidth/2 - 230
	y := screenHeight/2 - 200

	text.Draw(screen,
		"TANK",
		pReg.TitleFont,
		x, y,
		color.White)

	y += 100
	text.Draw(screen,
		fmt.Sprintf("LocalAddr %v", pReg.LocalAddr),
		pReg.ArcadeFont,
		x, y,
		color.White)

	y += 50
	nameY := y
	nameText := fmt.Sprintf("Name: %v", pReg.Name)
	text.Draw(screen,
		nameText,
		pReg.ArcadeFont,
		x, y,
		color.White)

	y += 100
	text.Draw(screen,
		"Reg",
		pReg.ArcadeFont,
		x, y,
		color.White)

	y += 50
	ipY := y
	ipText := fmt.Sprintf("IP  : %v", pReg.Ip)
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

		switch pReg.InputFlag {
		case InputFlagName:
			t = nameText
			y = nameY
		case InputFlagIp:
			t = ipText
			y = ipY
		case InputFlagPort:
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
