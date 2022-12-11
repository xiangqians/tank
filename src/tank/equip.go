// equip
// @author xiangqian
// @date 23:10 2022/12/10
package tank

import (
	"github.com/hajimehoshi/ebiten"
	"log"
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
const MaxEquipCount uint8 = 255

type Equip struct {
	*AbsGraphics
	EquipType     EquipType `json:"equipType"` // 装备类型
	EffectiveTime uint8     `json:"timeout"`   // 拾取装备后有效时间，单位s
}

func CreateEquip() *Equip {
	x, y := RandXY()
	pEquip := &Equip{
		AbsGraphics:   CreateAbsGraphics(pApp.pReg.Name, GraphicsTyEquip, Location{float64(x), float64(y)}, DirectionRight, DefaultTankSpeed),
		EquipType:     RandEquipType(),
		EffectiveTime: 10,
	}
	pEquip.Init(pEquip)
	return pEquip
}

func RandEquipType() EquipType {
	return EquipTypeBulletAcc
}

func EquipGenerator() {
	//time.Sleep(10 * time.Second)
	for {
		//time.Sleep(time.Duration(RandIntn(60)) * time.Second)
		//time.Sleep(time.Duration(RandIntn(6)) * time.Second)
		time.Sleep(time.Duration(10) * time.Millisecond)
		if pApp.pGame.EquipCount < MaxEquipCount {
			pEquip := CreateEquip()

			// 通知其它端点
			pApp.pEndpoint.SendGraphicsToAddrs(pEquip)

			pApp.pGame.AddGraphics(pEquip)
		}
	}
}

// 佩戴装备
func (pEquip *Equip) WearEquip(pTank *Tank) {
	if pTank.EquipMap == nil {
		log.Printf("初始化EquipMap\n")
		pTank.EquipMap = make(map[string]*Equip, 8)
	}

	//if v, r := pTank.EquipMap[pEquip.Id]; r {
	if _, r := pTank.EquipMap[pEquip.Id]; r {
		return
	}

	pTank.EquipMap[pEquip.Id] = pEquip

	switch pEquip.EquipType {
	// 坦克加速
	case EquipTypeTankAcc:

	// 坦克减速
	case EquipTypeTankDec:
		// 子弹加速
	case EquipTypeBulletAcc:

		// 子弹减速
	case EquipTypeBulletDec:
	}

	log.Printf("WearEquip: %v\n", pEquip.Id)

}

// 卸下装备
func (pEquip *Equip) RemoveEquip(pTank *Tank) {
	log.Printf("RemoveEquip: %v\n", pEquip.Id)
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
