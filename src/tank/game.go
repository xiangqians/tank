// Game Home
// @author xiangqian
// @date 13:12 2022/12/03
package tank

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"time"
)

type Game struct {
	GraphicsMap     map[string]Graphics      // 图形Map
	GraphicsMapChan chan map[string]Graphics // 图形Map channel
	pTank           *Tank                    // 当前用户坦克
}

func (pGame *Game) Init() {
	pGame.GraphicsMap = make(map[string]Graphics, 8)
	pGame.GraphicsMapChan = make(chan map[string]Graphics, 1)
	pGame.GraphicsMapChan <- pGame.GraphicsMap

	// TANK
	pGame.pTank = CreateTank(Location{300, 100}, DirectionRight, SpeedNormal)
	pGame.AddGraphics(pGame.pTank)

	// test
	//pGame.AddGraphics(CreateTank(Location{400, 200}, DirectionRight, SpeedNormal))

	// 清理
	go func() {
		for {
			pGame.Clean()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (pGame *Game) Clean() {
	// 阻塞获取 chanel 中的 map
	graphicsMap := <-pGame.GraphicsMapChan

	// 再将 map 添加到 channel
	defer func() { pGame.GraphicsMapChan <- graphicsMap }()

	// 定义id切片（slice）
	var ids []string
	ids = nil
	index := 0
	for id, graphics := range graphicsMap {
		if graphics.GetStatus() == StatusTerm {
			if ids == nil {
				ids = make([]string, len(graphicsMap))
			}
			ids[index] = id
			index++
		}
	}

	if ids != nil {
		for i := 0; i < index; i++ {
			id := ids[i]
			delete(graphicsMap, id)
			//fmt.Printf("delete %v\n", id)
		}
	}
}

func (pGame *Game) AddAbsGraphics(pAbsGraphics *AbsGraphics) {
	var graphics Graphics = nil
	switch pAbsGraphics.GraphicsTy {
	case GraphicsTyTank:
		pTank := &Tank{AbsGraphics: pAbsGraphics}
		pTank.Init(pTank)
		graphics = pTank

	case GraphicsTyBullet:
		pBullet := &Bullet{AbsGraphics: pAbsGraphics}
		pBullet.Init(pBullet)
		graphics = pBullet
	}

	if graphics == nil {
		return
	}

	// 如果是当前坦克时
	if graphics.GetId() == pGame.pTank.GetId() {
		pGame.pTank.Status = graphics.GetStatus()
		pGame.pTank.Hp = graphics.GetHp()
		return
	}

	// add
	pGame.AddGraphics(graphics)
}

func (pGame *Game) AddGraphics(graphics Graphics) {
	// 阻塞获取 chanel 中的 map
	graphicsMap := <-pGame.GraphicsMapChan

	// 再将 map 添加到 channel
	defer func() { pGame.GraphicsMapChan <- graphicsMap }()

	// 删除 map key
	delete(graphicsMap, graphics.GetId())
	// 添加 map key
	graphicsMap[graphics.GetId()] = graphics
}

func (pGame *Game) DelGraphics(graphics Graphics) {
	// 阻塞获取 chanel 中的 map
	graphicsMap := <-pGame.GraphicsMapChan

	// 再将 map 添加到 channel
	defer func() { pGame.GraphicsMapChan <- graphicsMap }()

	// 删除 map key
	delete(graphicsMap, graphics.GetId())
}

func (pGame *Game) Update(screen *ebiten.Image) error {
	if pGame.pTank.Status == StatusTerm {
		return nil
	}

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
