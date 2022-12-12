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
	EquipCount      uint8                    // 装备数量
}

// 端点默认监听端口
var DefaultEndpointPort int = 0

// 坦克默认速度
var DefaultTankSpeed = SpeedNormal

// 子弹默认速度
var DefaultBulletSpeed = SpeedNormal

func (pGame *Game) Init() {
	pGame.GraphicsMap = make(map[string]Graphics, 8)
	pGame.GraphicsMapChan = make(chan map[string]Graphics, 1)
	pGame.GraphicsMapChan <- pGame.GraphicsMap

	// TANK
	pGame.pTank = CreateDefaultTank()
	pGame.AddGraphics0(pGame.pTank)

	// test
	//pGame.AddGraphics0(CreateEquip())

	// 装备生成器
	go EquipGenerator()

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

func (pGame *Game) AddGraphics(graphics Graphics) {
	// 如果是当前坦克时
	if graphics.GetId() == pGame.pTank.GetId() {
		pGame.pTank.Status = graphics.GetStatus()
		pGame.pTank.Hp = graphics.GetHp()
		return
	}

	pGame.AddGraphics0(graphics)
}

func (pGame *Game) AddGraphics0(graphics Graphics) {
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
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			pGame.pTank.Reset()
			pGame.AddGraphics(pGame.pTank)
			pApp.pEndpoint.SendGraphicsToAddrs(pGame.pTank)
		}
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
	//if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
	// 支持长按空格键发射子弹
	if pApp.IsSpaceKeyPressed() {
		//log.Printf("space\n")
		pGame.pTank.Fire()
	}

	return nil
}

func (pGame *Game) Draw(screen *ebiten.Image) {
	pGame._Draw(screen)
	pGame.Draw1(screen)
}

// 游戏界面闪烁解决了，但出现另一个问题：图形渲染变慢了。
func (pGame *Game) Draw1(screen *ebiten.Image) {
	// 阻塞获取 chanel 中的 map
	graphicsMap := <-pGame.GraphicsMapChan

	// 再将 map 添加到 channel
	defer func() { pGame.GraphicsMapChan <- graphicsMap }()

	var equipCount uint8 = 0
	for _, graphics := range graphicsMap {
		graphics.VerifyTimestamp()
		if graphics.GetStatus() != StatusTerm {
			if graphics.GetGraphicsTy() == GraphicsTyEquip {
				equipCount++
			}
			graphics.Draw(screen)
		}
	}
	pGame.EquipCount = equipCount
}

// 会出现游戏界面闪烁问题
func (pGame *Game) Draw0(screen *ebiten.Image) {
	select {
	// 非阻塞获取 chanel 中的 map
	case graphicsMap := <-pGame.GraphicsMapChan:
		// 再将 map 添加到 channel
		defer func() { pGame.GraphicsMapChan <- graphicsMap }()
		var equipCount uint8 = 0
		for _, graphics := range graphicsMap {
			graphics.VerifyTimestamp()
			if graphics.GetStatus() != StatusTerm {
				if graphics.GetGraphicsTy() == GraphicsTyEquip {
					equipCount++
				}
				graphics.Draw(screen)
			}
		}
		pGame.EquipCount = equipCount
	default:
	}
}

func (pGame *Game) _Draw(screen *ebiten.Image) {
	// 坦克速度描述
	ts, _ := SpeedString(DefaultTankSpeed)
	if pGame.pTank.Speed != DefaultTankSpeed {
		str, _ := SpeedString(pGame.pTank.Speed)
		ts = fmt.Sprintf("%v -> %v (%v s)", ts, str, int64(pGame.pTank.TankAccEquipEffectiveTime)-(time.Now().UnixNano()-pGame.pTank.TankAccEquipGetTimestamp)/int64(time.Second))
	}

	// 子弹速度描述
	bs, _ := SpeedString(DefaultBulletSpeed)
	if pGame.pTank.BulletSpeed != DefaultBulletSpeed {
		str, _ := SpeedString(pGame.pTank.BulletSpeed)
		bs = fmt.Sprintf("%v -> %v (%v s)", bs, str, int64(pGame.pTank.BulletAccEquipEffectiveTime)-(time.Now().UnixNano()-pGame.pTank.BulletAccEquipGetTimestamp)/int64(time.Second))
	}

	// 坦克隐形描述
	ti := fmt.Sprintf("%v", pGame.pTank.TankInvisFlag)
	if pGame.pTank.TankInvisFlag {
		ti = fmt.Sprintf("%v (%v s)", pGame.pTank.TankInvisFlag, int64(pGame.pTank.TankInvisEquipEffectiveTime)-(time.Now().UnixNano()-pGame.pTank.TankInvisEquipGetTimestamp)/int64(time.Second))
	}

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(""+
			"LocalAddr : %v\n"+
			"Name      : %v\n"+
			"TS        : %v\n"+
			"BS        : %v\n"+
			"TI        : %v\n"+
			"HP        : %v\n",
			pApp.pReg.LocalAddr,
			pGame.pTank.GetName(),
			ts,
			bs,
			ti,
			pGame.pTank.GetHp()))
}
