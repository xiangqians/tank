// application
// @author xiangqian
// @date 22:22 2022/11/29
package tank

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"log"
	"time"
)

const (
	// 1280 * 720
	screenWidth  = 1280
	screenHeight = 720
)

// 端点
var endpoint *Endpoint

// 图形Map
var graphicsMap map[string]Graphics

// 图形Map channel
var graphicsMapChan chan map[string]Graphics

// 当前用户坦克
var pTank *Tank

func init() {
	endpoint = &Endpoint{}

	graphicsMap = make(map[string]Graphics, 8)
	graphicsMapChan = make(chan map[string]Graphics, 1)
	graphicsMapChan <- graphicsMap
}

func addGraphics(graphics Graphics) {
	graphicsMap[graphics.Id()] = graphics
}

type Game struct {
	prevKeyPressedUnixNano int64
	curKeyPressedUnixNano  int64
}

func (game *Game) Update(screen *ebiten.Image) error {

	// up
	if game.IsKeyPressed(ebiten.KeyUp) {
		//log.Printf("up\n")
		pTank.Move(DirectionUp)

	} else
	// down
	if game.IsKeyPressed(ebiten.KeyDown) {
		//log.Printf("down\n")
		pTank.Move(DirectionDown)

	} else
	// left
	if game.IsKeyPressed(ebiten.KeyLeft) {
		//log.Printf("left\n")
		pTank.Move(DirectionLeft)

	} else
	// right
	if game.IsKeyPressed(ebiten.KeyRight) {
		//log.Printf("right\n")
		pTank.Move(DirectionRight)

	}

	// space
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		//log.Printf("space\n")
		pTank.Fire()
	}

	return nil
}

func (game *Game) IsKeyPressed(key ebiten.Key) bool {
	if ebiten.IsKeyPressed(key) {
		game.curKeyPressedUnixNano = time.Now().UnixNano()
		result := game.curKeyPressedUnixNano-game.prevKeyPressedUnixNano >= 10
		game.prevKeyPressedUnixNano = game.curKeyPressedUnixNano
		return result
	}
	return false
}

func (game *Game) Draw(screen *ebiten.Image) {
	select {
	// 非阻塞获取 chanel 中的 map
	case graphicsMap := <-graphicsMapChan:
		// 再将 map 添加到 channel
		defer func() { graphicsMapChan <- graphicsMap }()
		for _, graphics := range graphicsMap {
			if graphics.Status() != StatusTerm {
				graphics.Draw(screen)
			}
		}

	default:
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Run() {
	// 端点监听
	go endpoint.Listen()

	pTank = CreateTank(Location{300, 100}, DirectionRight, SpeedNormal)
	addGraphics(pTank)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tank")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
