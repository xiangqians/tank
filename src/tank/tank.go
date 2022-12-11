// tank
// 坦克
//
// @author xiangqian
// @date 10:58 2022/11/27
package tank

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"time"
)

type Tank struct {
	*AbsGraphics
	BulletSpeed Speed // 子弹速度
	// 坦克装备之坦克加速
	TankAccEquipGetTimestamp  int64 // 拾起装备时间戳，单位，ns
	TankAccEquipEffectiveTime uint8 // 拾起装备后有效时间，单位s

	// 子弹装备之子弹加速
	BulletAccEquipGetTimestamp  int64 // 拾起装备时间戳，单位，ns
	BulletAccEquipEffectiveTime uint8 // 拾起装备后有效时间，单位s
	//pTankInvisEquip *Equip // 坦克隐形装备
}

func CreateDefaultTank() *Tank {
	x, y := RandXY()
	pTank := &Tank{
		AbsGraphics:                 CreateAbsGraphics(pApp.pReg.Name, GraphicsTyTank, Location{float64(x), float64(y)}, DirectionUp, DefaultTankSpeed),
		BulletSpeed:                 DefaultBulletSpeed,
		TankAccEquipEffectiveTime:   10,
		BulletAccEquipEffectiveTime: 10,
	}
	pTank.Init()

	go func() {
		for {
			pTank.VerifyEquip()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return pTank
}

func (pTank *Tank) Init() {
	pTank.AbsGraphics.Init(pTank)
}

// RemoveEquip
func (pTank *Tank) VerifyEquip() {
	// 坦克装备之坦克加速
	if pTank.Speed != DefaultTankSpeed && time.Now().UnixNano()-pTank.TankAccEquipGetTimestamp > int64(time.Duration(pTank.TankAccEquipEffectiveTime)*time.Second) {
		pTank.Speed = DefaultTankSpeed

		// 卸下装备
		log.Printf("RemoveEquip: TankAcc\n")
	}

	// 子弹装备之子弹加速
	if pTank.BulletSpeed != DefaultBulletSpeed && time.Now().UnixNano()-pTank.BulletAccEquipGetTimestamp > int64(time.Duration(pTank.BulletAccEquipEffectiveTime)*time.Second) {
		pTank.BulletSpeed = DefaultBulletSpeed

		// 卸下装备
		log.Printf("RemoveEquip: BulletAcc\n")
	}
}

// 佩戴装备
func (pTank *Tank) WearEquip(pEquip *Equip) {
	switch pEquip.EquipType {
	// 坦克加速
	case EquipTypeTankAcc:
		pTank.TankAccEquipGetTimestamp = time.Now().UnixNano()
		pTank.Speed = SpeedFast
		log.Printf("WearEquip: TankAcc\n")

	// 子弹加速
	case EquipTypeBulletAcc:
		pTank.BulletAccEquipGetTimestamp = time.Now().UnixNano()
		pTank.BulletSpeed = SpeedFast
		log.Printf("WearEquip: BulletAcc\n")
	}
}

// 重置（恢复）坦克
func (pTank *Tank) Reset() {
	pTank.Hp = DefaultHp
	x, y := RandXY()
	pTank.Location = Location{float64(x), float64(y)}
	pTank.Status = StatusRun
}

// 开火
func (pTank *Tank) Fire() {
	pBullet := CreateBullet(pTank, pTank.BulletSpeed)
	pApp.pGame.AddGraphics(pBullet)
	go pBullet.Run()
}

func (pTank *Tank) Move(direction Direction) {
	pTank.AbsGraphics.Move(direction)
	pApp.pEndpoint.SendGraphicsToAddrs(pTank)
}

func (pTank *Tank) UpImage() *ebiten.Image {
	return pApp.pImage.pTankUpImage
}

func (pTank *Tank) DownImage() *ebiten.Image {
	return pApp.pImage.pTankDownImage
}

func (pTank *Tank) LeftImage() *ebiten.Image {
	return pApp.pImage.pTankLeftImage
}

func (pTank *Tank) RightImage() *ebiten.Image {
	return pApp.pImage.pTankRightImage
}
