// Game Home
// @author xiangqian
// @date 13:12 2022/12/03
package tank

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"sync"
	"time"
)

// 图形仓库状态
type GRepoStatus int8

const (
	GRepoStatusReading GRepoStatus = iota + 1 // 读取中
	GRepoStatusRead                           // 已读
	GRepoStatusWriting                        // 写入中
	GRepoStatusWritten                        // 已写
)

// 图形仓库
type GRepo struct {
	Arr         []Graphics  // 图形切片
	Len         int         // 图形切片长度
	GRepoStatus GRepoStatus // 图形仓库状态
}

type Game struct {
	GraphicsMap     map[string]Graphics      // 图形Map
	GraphicsMapChan chan map[string]Graphics // 图形Map channel

	// 图形仓库切片
	pGRepoArr    []*GRepo   // 图形仓库切片
	GRepoArrLock sync.Mutex // 图形仓库切片互斥锁

	pTank      *Tank  // 当前用户坦克
	EquipCount uint16 // 装备数量
}

// 端点默认监听端口
var DefaultEndpointPort int = 0

// 坦克默认速度
var DefaultTankSpeed = SpeedNormal

// 子弹默认速度
var DefaultBulletSpeed = SpeedNormal

func (pGame *Game) Init() {
	pGame.GraphicsMap = make(map[string]Graphics, 64)
	pGame.GraphicsMapChan = make(chan map[string]Graphics, 1)
	pGame.GraphicsMapChan <- pGame.GraphicsMap
	//pGame.pGRepoArr = []*GRepo{
	//	{
	//		Arr:         make([]Graphics, 64),
	//		Len:         0,
	//		GRepoStatus: GRepoStatusRead,
	//	},
	//	{
	//		Arr:         make([]Graphics, 64),
	//		Len:         0,
	//		GRepoStatus: GRepoStatusRead,
	//	},
	//	{
	//		Arr:         make([]Graphics, 64),
	//		Len:         0,
	//		GRepoStatus: GRepoStatusRead,
	//	},
	//	{
	//		Arr:         make([]Graphics, 64),
	//		Len:         0,
	//		GRepoStatus: GRepoStatusRead,
	//	},
	//	{
	//		Arr:         make([]Graphics, 64),
	//		Len:         0,
	//		GRepoStatus: GRepoStatusRead,
	//	},
	//	{
	//		Arr:         make([]Graphics, 64),
	//		Len:         0,
	//		GRepoStatus: GRepoStatusRead,
	//	},
	//}

	// TANK
	pGame.pTank = CreateDefaultTank()
	pGame.AddGraphics(pGame.pTank)

	// test
	//pGame.AddGraphics(CreateEquip())

	// 写入图形仓库
	//go func() {
	//	for {
	//		pGame.WriteGRepo()
	//		time.Sleep(10 * time.Millisecond)
	//	}
	//}()

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

// 获取可读的图形仓库
func (pGame *Game) GetReadableGRepo() *GRepo {
	// 加锁
	pGame.GRepoArrLock.Lock()

	// 解锁
	defer pGame.GRepoArrLock.Unlock()

	var pReadableGRepo *GRepo = nil
	var pReadableGRepoBak *GRepo = nil
	for _, pGRepo := range pGame.pGRepoArr {
		if pGRepo.GRepoStatus == GRepoStatusWritten {
			pReadableGRepo = pGRepo
			break
		}
		if pGRepo.GRepoStatus == GRepoStatusRead {
			pReadableGRepoBak = pGRepo
		}
	}

	if pReadableGRepo == nil {
		pReadableGRepo = pReadableGRepoBak
	}

	if pReadableGRepo != nil {
		pReadableGRepo.GRepoStatus = GRepoStatusReading
	}

	return pReadableGRepo
}

// 获取可写的图形仓库
func (pGame *Game) GetWritableGRepo() *GRepo {
	// 加锁
	pGame.GRepoArrLock.Lock()

	// 解锁
	defer pGame.GRepoArrLock.Unlock()

	var pWritableGRepo *GRepo = nil
	for _, pGRepo := range pGame.pGRepoArr {
		if pGRepo.GRepoStatus == GRepoStatusRead {
			pWritableGRepo = pGRepo
			break
		}
	}

	if pWritableGRepo != nil {
		pWritableGRepo.GRepoStatus = GRepoStatusWriting
	}

	return pWritableGRepo
}

func (pGame *Game) WriteGRepo() {
	select {
	// 非阻塞获取 chanel 中的 map
	case graphicsMap := <-pGame.GraphicsMapChan:
		pWritableGRepo := pGame.GetWritableGRepo()

		// 再将 map 添加到 channel
		defer func() {
			pGame.GraphicsMapChan <- graphicsMap
			if pWritableGRepo != nil {
				pWritableGRepo.GRepoStatus = GRepoStatusWritten
			}
		}()

		if pWritableGRepo != nil {
			i := 0
			l := len(pWritableGRepo.Arr)
			for _, graphics := range graphicsMap {
				if i < l {
					pWritableGRepo.Arr[i] = graphics
				} else {
					pWritableGRepo.Arr = append(pWritableGRepo.Arr, graphics)
				}
				i++
			}
			pWritableGRepo.Len = i
		}
	default:
	}
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
		pTank := graphics.(*Tank)
		pGame.pTank.Status = graphics.GetStatus()
		pGame.pTank.Hp = graphics.GetHp()
		pGame.pTank.BulletSpeed = pTank.BulletSpeed
		pGame.pTank.TankInvisFlag = pTank.TankInvisFlag
		graphics = pGame.pTank
	}

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

		//str := ""
		//for _, pGRepo := range pGame.pGRepoArr {
		//	if str != "" {
		//		str += ","
		//	}
		//	str += strconv.Itoa(int(pGRepo.GRepoStatus))
		//}
		//log.Printf("GRepoArr: %v\n", str)
	}

	return nil
}

func (pGame *Game) Draw(screen *ebiten.Image) {
	pGame._Draw(screen)
	pGame.Draw1(screen)
}

func (pGame *Game) Draw2(screen *ebiten.Image) {
	pReadableGRepo := pGame.GetReadableGRepo()
	if pReadableGRepo == nil {
		return
	}

	defer func() { pReadableGRepo.GRepoStatus = GRepoStatusRead }()

	var equipCount uint16 = 0
	for i, l := 0, pReadableGRepo.Len; i < l; i++ {
		graphics := pReadableGRepo.Arr[i]
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

// 游戏界面闪烁解决了，但出现另一个问题：图形渲染变慢了。
func (pGame *Game) Draw1(screen *ebiten.Image) {
	// 阻塞获取 chanel 中的 map
	graphicsMap := <-pGame.GraphicsMapChan

	// 再将 map 添加到 channel
	defer func() { pGame.GraphicsMapChan <- graphicsMap }()

	var equipCount uint16 = 0
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
		var equipCount uint16 = 0
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
