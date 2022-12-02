// Game
// @author xiangqian
// @date 22:22 2022/11/29
package tank

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// 游戏步骤
type GameStep int8

const (
	GameStepReg  GameStep = iota // 注册到端点界面
	GameStepHome                 // 游戏主界面
)

const (
	// 1280 * 720
	screenWidth  = 1280
	screenHeight = 720
)

var arcadeFont font.Face

// 端点
var pEndpoint *Endpoint

// 图形Map
var graphicsMap map[string]Graphics

// 图形Map channel
var graphicsMapChan chan map[string]Graphics

// 当前用户坦克
var pTank *Tank

var pGameReg *GameReg

func init() {
	logger()

	pGameReg = &GameReg{}

	pEndpoint = &Endpoint{}

	graphicsMap = make(map[string]Graphics, 8)
	graphicsMapChan = make(chan map[string]Graphics, 1)
	graphicsMapChan <- graphicsMap

	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func logger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pLogFile, err := os.OpenFile("./tank.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("open log file err, %v", err)
		return
	}

	log.SetOutput(pLogFile)
}

func addGraphics(graphics Graphics) {
	graphicsMap[graphics.GetId()] = graphics
}

type GameReg struct {
	LocalAddr    string
	A, B, C, D   byte
	Port         string
	CursorStatus bool
}

func (pGameReg *GameReg) Append(a string) {
	if len(pGameReg.Port) >= 5 {
		return
	}

	pGameReg.Port += a
}

type Game struct {
	GameStep               GameStep
	PrevKeyPressedUnixNano int64
	CurKeyPressedUnixNano  int64
}

func (pGame *Game) Init() {
	go func() {
		for pGame.GameStep == GameStepReg {
			pGameReg.CursorStatus = !pGameReg.CursorStatus
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func (pGame *Game) Update(screen *ebiten.Image) error {
	switch pGame.GameStep {
	case GameStepReg:
		pGame.UpdateReg(screen)
	case GameStepHome:
		pGame.UpdateHome(screen)
	}
	return nil
}

func (pGame *Game) UpdateReg(screen *ebiten.Image) {

	if inpututil.IsKeyJustPressed(ebiten.KeyKP0) || inpututil.IsKeyJustPressed(ebiten.Key0) {
		pGameReg.Append("0")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP1) || inpututil.IsKeyJustPressed(ebiten.Key1) {
		pGameReg.Append("1")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP2) || inpututil.IsKeyJustPressed(ebiten.Key2) {
		pGameReg.Append("2")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP3) || inpututil.IsKeyJustPressed(ebiten.Key3) {
		pGameReg.Append("3")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP4) || inpututil.IsKeyJustPressed(ebiten.Key4) {
		pGameReg.Append("4")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP5) || inpututil.IsKeyJustPressed(ebiten.Key5) {
		pGameReg.Append("5")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP6) || inpututil.IsKeyJustPressed(ebiten.Key6) {
		pGameReg.Append("6")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP7) || inpututil.IsKeyJustPressed(ebiten.Key7) {
		pGameReg.Append("7")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP8) || inpututil.IsKeyJustPressed(ebiten.Key8) {
		pGameReg.Append("8")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyKP9) || inpututil.IsKeyJustPressed(ebiten.Key9) {
		pGameReg.Append("9")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		length := len(pGameReg.Port)
		if length > 0 {
			pGameReg.Port = pGameReg.Port[:length-1]
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if len(pGameReg.Port) > 0 {
			pDgPkt := &DgPkt{}
			pDgPkt.DgPktTy = DgPktTyReg

			port, _ := strconv.ParseInt(pGameReg.Port, 10, 64)
			addr := &net.UDPAddr{
				IP:   pEndpoint.LocalAddr.IP,
				Port: int(port),
			}
			log.Printf("reg addr: %v\n", addr.String())
			pEndpoint.SendDgPkt(pDgPkt, addr)
			pGame.GameStep = GameStepHome
		}
	}
}

func (pGame *Game) UpdateHome(screen *ebiten.Image) {

	// up
	if pGame.IsKeyPressed(ebiten.KeyUp) {
		//log.Printf("up\n")
		pTank.Move(DirectionUp)

	} else
	// down
	if pGame.IsKeyPressed(ebiten.KeyDown) {
		//log.Printf("down\n")
		pTank.Move(DirectionDown)

	} else
	// left
	if pGame.IsKeyPressed(ebiten.KeyLeft) {
		//log.Printf("left\n")
		pTank.Move(DirectionLeft)

	} else
	// right
	if pGame.IsKeyPressed(ebiten.KeyRight) {
		//log.Printf("right\n")
		pTank.Move(DirectionRight)

	}

	// space
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		//log.Printf("space\n")
		pTank.Fire()
	}
}

func (pGame *Game) IsKeyPressed(key ebiten.Key) bool {
	if ebiten.IsKeyPressed(key) {
		pGame.CurKeyPressedUnixNano = time.Now().UnixNano()
		result := pGame.CurKeyPressedUnixNano-pGame.PrevKeyPressedUnixNano >= 10
		pGame.PrevKeyPressedUnixNano = pGame.CurKeyPressedUnixNano
		return result
	}
	return false
}

func (pGame *Game) Draw(screen *ebiten.Image) {
	switch pGame.GameStep {
	case GameStepReg:
		pGame.DrawReg(screen)
	case GameStepHome:
		pGame.DrawHome(screen)
	}
}

func (pGame *Game) DrawReg(screen *ebiten.Image) {

	x := screenWidth/2 - 230
	y := screenHeight/2 - 200

	text.Draw(screen,
		fmt.Sprintf("localAddr %v", pGameReg.LocalAddr),
		arcadeFont,
		x, y,
		color.White)

	y += 100
	text.Draw(screen,
		fmt.Sprintf("IP:   %v", strings.Split(pGameReg.LocalAddr, ":")[0]),
		arcadeFont,
		x, y,
		color.White)

	y += 50
	t := fmt.Sprintf("Port: %v", pGameReg.Port)
	text.Draw(screen,
		t,
		arcadeFont,
		x, y,
		color.White)

	if pGameReg.CursorStatus {
		text.Draw(screen,
			fmt.Sprintf("|"),
			arcadeFont,
			x+len(t)*16, y,
			color.White)
	}
}

func (pGame *Game) DrawHome(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("localAddr %v", pGameReg.LocalAddr))

	select {
	// 非阻塞获取 chanel 中的 map
	case graphicsMap := <-graphicsMapChan:
		// 再将 map 添加到 channel
		defer func() { graphicsMapChan <- graphicsMap }()
		for _, graphics := range graphicsMap {
			if graphics.GetStatus() != StatusTerm {
				graphics.Draw(screen)
			}
		}
	default:
	}
}

func (pGame *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Run() {
	// 端点监听
	go pEndpoint.Listen()

	pTank = CreateTank(Location{300, 100}, DirectionRight, SpeedNormal)
	addGraphics(pTank)

	pGame := &Game{}
	pGame.Init()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tank")
	if err := ebiten.RunGame(pGame); err != nil {
		log.Fatal(err)
	}
}
