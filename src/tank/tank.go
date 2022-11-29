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
	pTank.sub = pTank
	return pTank
}

// 开火
func (tank *Tank) Fire() {
}

func (tank *Tank) UpImg() *ebiten.Image {
	return pTankUpImg
}

func (tank *Tank) DownImg() *ebiten.Image {
	return pTankDownImg
}

func (tank *Tank) LeftImg() *ebiten.Image {
	return pTankLeftImg
}

func (tank *Tank) RightImg() *ebiten.Image {
	return pTankRightImg
}
