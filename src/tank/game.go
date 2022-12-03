// Game Home
// @author xiangqian
// @date 13:12 2022/12/03
package tank

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Game struct {
	GraphicsMap     map[string]Graphics      // 图形Map
	GraphicsMapChan chan map[string]Graphics // 图形Map channel
	pTank           *Tank                    // 当前用户坦克
}

func (pGame *Game) init() {
	pGame.GraphicsMap = make(map[string]Graphics, 8)
	pGame.GraphicsMapChan = make(chan map[string]Graphics, 1)
	pGame.GraphicsMapChan <- pGame.GraphicsMap

	// TANK
	pGame.pTank = CreateTank(Location{300, 100}, DirectionRight, SpeedNormal)
	pGame.AddGraphics(pGame.pTank)
}

func (pGame *Game) AddGraphics(graphics Graphics) {
	pGame.GraphicsMap[graphics.GetId()] = graphics
}

func (pGame *Game) Update(screen *ebiten.Image) error {
	// up
	if pApp.IsKeyPressed(ebiten.KeyUp) {
		//log.Printf("up\n")
		pGame.pTank.Move(DirectionUp)

	} else
	// down
	if pApp.IsKeyPressed(ebiten.KeyDown) {
		//log.Printf("down\n")
		pGame.pTank.Move(DirectionDown)

	} else
	// left
	if pApp.IsKeyPressed(ebiten.KeyLeft) {
		//log.Printf("left\n")
		pGame.pTank.Move(DirectionLeft)

	} else
	// right
	if pApp.IsKeyPressed(ebiten.KeyRight) {
		//log.Printf("right\n")
		pGame.pTank.Move(DirectionRight)
	}

	// space
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		//log.Printf("space\n")
		pGame.pTank.Fire()
	}

	return nil
}

func (pGame *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("localAddr %v\nHP: %v", pApp.pReg.LocalAddr, pGame.pTank.GetHp()))

	select {
	// 非阻塞获取 chanel 中的 map
	case graphicsMap := <-pGame.GraphicsMapChan:
		// 再将 map 添加到 channel
		defer func() { pGame.GraphicsMapChan <- graphicsMap }()
		for _, graphics := range graphicsMap {
			if graphics.GetStatus() != StatusTerm {
				graphics.Draw(screen)
			}
		}
	default:
	}
}
