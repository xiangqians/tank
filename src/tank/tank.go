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
	pTank := &Tank{AbsGraphics: CreateAbsGraphics("3rwdc2335cv", GraphicsTyTank, location, direction, speed)}
	pTank.Init(pTank)
	return pTank
}

// 开火
func (pTank *Tank) Fire() {
	pBullet := CreateBullet(pTank, SpeedNormal)
	pApp.pGame.AddGraphics(pBullet)
	go pBullet.Run()
}

func (pTank *Tank) Move(direction Direction) {
	pTank.AbsGraphics.Move(direction)
	pApp.pEndpoint.SendGraphics(pTank)
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
