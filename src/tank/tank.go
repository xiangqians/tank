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
	pLocation *Location
	pOptions  *ebiten.DrawImageOptions
	pImg      *ebiten.Image
}

func CreateTank(location Location, direction Direction, speed Speed) *Tank {
	pTank := &Tank{
		AbsGraphics: &AbsGraphics{
			id:        "tank_" + Uuid(),
			direction: direction,
			speed:     speed,
			status:    StatusNew,
		},
		pLocation: &Location{location.x, location.y},
		pOptions:  nil,
		pImg:      nil,
	}
	return pTank
}

func (tank *Tank) Move(direction Direction) {
	// -→ x
	// ↓ y
	location := tank.pLocation
	tank.direction = direction
	var pImg *ebiten.Image
	var xx float64 = 1 + float64(tank.speed)
	switch direction {
	case DirectionUp:
		location.y -= xx
		pImg = pTankUpImg

	case DirectionDown:
		location.y += xx
		pImg = pTankDownImg

	case DirectionLeft:
		location.x -= xx
		pImg = pTankLeftImg

	case DirectionRight:
		location.x += xx
		pImg = pTankRightImg
	}
	tank.pImg = pImg
}

// 开火
func (tank *Tank) Fire() {
}

func (tank *Tank) Draw(screen *ebiten.Image) error {
	if err := tank.AbsGraphics.Draw(screen); err != nil {
		return err
	}
	options := tank.pOptions
	if options == nil {
		options = &ebiten.DrawImageOptions{}
		tank.pOptions = options
	}
	options.GeoM.Reset()
	location := tank.pLocation
	options.GeoM.Translate(location.x, location.y)

	pImg := tank.pImg
	if pImg == nil {
		switch tank.direction {
		case DirectionUp:
			pImg = pTankUpImg

		case DirectionDown:
			pImg = pTankDownImg

		case DirectionLeft:
			pImg = pTankLeftImg

		case DirectionRight:
			pImg = pTankRightImg
		}
		tank.pImg = pImg
	}
	screen.DrawImage(pImg, options)
	return nil
}
