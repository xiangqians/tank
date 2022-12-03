// Reg
// @author xiangqian
// @date 13:12 2022/12/03
package tank

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
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

type Reg struct {
	LocalAddr    string
	A, B, C, D   byte
	Port         string
	ArcadeFont   font.Face
	CursorStatus bool
}

func (pReg *Reg) Init() {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	pReg.ArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	go func() {
		for pApp.appStep == AppStepReg {
			pReg.CursorStatus = !pReg.CursorStatus
			time.Sleep(500 * time.Millisecond)
		}
	}()
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
	} else if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		length := len(pReg.Port)
		if length > 0 {
			pReg.Port = pReg.Port[:length-1]
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if len(pReg.Port) > 0 {
			pDgPkt := &DgPkt{}
			pDgPkt.DgPktTy = DgPktTyReg

			port, _ := strconv.ParseInt(pReg.Port, 10, 64)
			addr := &net.UDPAddr{
				IP:   pApp.pEndpoint.pLocalAddr.IP,
				Port: int(port),
			}
			log.Printf("reg addr: %v\n", addr.String())
			pApp.pEndpoint.SendDgPkt(pDgPkt, addr)
			pApp.appStep = AppStepGame
		}
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
	text.Draw(screen,
		fmt.Sprintf("IP:   %v", strings.Split(pReg.LocalAddr, ":")[0]),
		pReg.ArcadeFont,
		x, y,
		color.White)

	y += 50
	t := fmt.Sprintf("Port: %v", pReg.Port)
	text.Draw(screen,
		t,
		pReg.ArcadeFont,
		x, y,
		color.White)

	if pReg.CursorStatus {
		text.Draw(screen,
			fmt.Sprintf("|"),
			pReg.ArcadeFont,
			x+len(t)*16, y,
			color.White)
	}
}

func (pReg *Reg) Append(a string) {
	if len(pReg.Port) >= 5 {
		return
	}

	pReg.Port += a
}
