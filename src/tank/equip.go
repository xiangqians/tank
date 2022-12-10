// equip
// @author xiangqian
// @date 23:10 2022/12/10
package tank

import (
	"github.com/hajimehoshi/ebiten"
	"time"
)

// 装备类型
type EquipType int8

const (
	EquipTypeTankAcc   EquipType = iota + 1 // 坦克加速
	EquipTypeTankDec                        // 坦克减速
	EquipTypeBulletAcc                      // 子弹加速
	EquipTypeBulletDec                      // 子弹减速
	EquipTypeHpInc                          // HP增加
	EquipTypeHpDec                          // HP减少
	EquipTypeTankInvis                      // 坦克隐形
)

// 装备最大数量
const MaxEquipCount uint8 = 5

type Equip struct {
	*AbsGraphics
	EquipType EquipType `json:"equipType"` // 装备类型
}

func CreateEquip() *Equip {
	x, y := RandXY()
	pEquip := &Equip{
		AbsGraphics: CreateAbsGraphics(pApp.pReg.Name, GraphicsTyEquip, Location{float64(x), float64(y)}, DirectionRight, DefaultTankSpeed),
		EquipType:   RandEquipType(),
	}
	pEquip.Init(pEquip)
	return pEquip
}

func RandEquipType() EquipType {
	return EquipTypeBulletAcc
}

func EquipGenerator() {
	time.Sleep(10 * time.Second)
	for {
		time.Sleep(time.Duration(RandIntn(60)) * time.Second)
		if pApp.pGame.EquipCount < MaxEquipCount {
			pEquip := CreateEquip()

			// 通知其它端点
			pApp.pEndpoint.SendGraphicsToAddrs(pEquip)

			pApp.pGame.AddGraphics(pEquip)
		}
	}
}

func (pEquip *Equip) Intersect(x, y float64, otherGraphics Graphics) bool {
	// 装备不判断图形之间是否相交
	return false
}

func (pEquip *Equip) UpImage() *ebiten.Image {
	return pApp.pImage.pEquipImage
}

func (pEquip *Equip) DownImage() *ebiten.Image {
	return pApp.pImage.pEquipImage
}

func (pEquip *Equip) LeftImage() *ebiten.Image {
	return pApp.pImage.pEquipImage
}

func (pEquip *Equip) RightImage() *ebiten.Image {
	return pApp.pImage.pEquipImage
}
