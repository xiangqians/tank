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
	EquipTypeTankInvis                      // 坦克隐形
	EquipTypeBulletAcc                      // 子弹加速
	EquipTypeHpInc                          // HP增加
)

var equipTypeIndex = 0
var equipTypes = []EquipType{EquipTypeTankAcc, EquipTypeBulletAcc, EquipTypeHpInc, EquipTypeTankInvis}

// 装备最大数量
// const MaxEquipCount uint8 = 5
const MaxEquipCount uint16 = 5

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
	pEquip.Init()
	return pEquip
}

func RandEquipType() EquipType {
	if equipTypeIndex >= len(equipTypes) {
		equipTypeIndex = 0
	}
	equipType := equipTypes[equipTypeIndex]
	equipTypeIndex++
	return equipType
}

func EquipGenerator() {
	for {
		time.Sleep(time.Duration(RandIntn(10)) * time.Second)
		//time.Sleep(50 * time.Millisecond)
		if pApp.pGame.EquipCount < MaxEquipCount {
			pEquip := CreateEquip()

			// 通知其它端点
			pApp.pEndpoint.SendGraphicsToAddrs(pEquip)

			pApp.pGame.AddGraphics(pEquip)
		}
	}
}

func (pEquip *Equip) Init() {
	pEquip.AbsGraphics.Init(pEquip)
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
