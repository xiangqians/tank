// tank
// 坦克
//
// @author xiangqian
// @date 10:58 2022/11/27
package tank

import (
	"github.com/hajimehoshi/ebiten"
)

type Tank struct {
	*AbsGraphics
}

func CreateTank(location Location, direction Direction, speed Speed) *Tank {
	pTank := &Tank{AbsGraphics: CreateAbsGraphics(location, direction, speed)}
	pTank.Init(pTank)
	return pTank
}

// 开火
func (pTank *Tank) Fire() {
	addGraphics(CreateBullet(pTank, SpeedNormal))
}

func (pTank *Tank) UpImg() *ebiten.Image {
	return pTankUpImg
}

func (pTank *Tank) DownImg() *ebiten.Image {
	return pTankDownImg
}

func (pTank *Tank) LeftImg() *ebiten.Image {
	return pTankLeftImg
}

func (pTank *Tank) RightImg() *ebiten.Image {
	return pTankRightImg
}