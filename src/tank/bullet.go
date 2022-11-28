// bullet
// 子弹
//
// @author xiangqian
// @date 11:01 2022/11/27
package tank

import (
	"github.com/nsf/termbox-go"
)

type Bullet struct {
	*AbsGraphics
	location *Location // 子弹位置
}

func CreateBullet(location Location, direction Direction, speed Speed) *Bullet {
	return &Bullet{
		AbsGraphics: &AbsGraphics{
			id:        "bullet_" + Uuid(),
			direction: direction,
			speed:     speed,
			alive:     true,
		},
		location: &Location{location.x, location.y},
	}
}

func (bullet *Bullet) Draw() error {

	if err := bullet.AbsGraphics.Draw(); err != nil {
		return err
	}

	location := bullet.location
	termbox.SetCell(location.x, location.y, '⬛', termbox.ColorRed, termbox.ColorRed)
	switch bullet.direction {
	case DirectionUp:
		location.y -= 1

	case DirectionDown:
		location.y += 1

	case DirectionLeft:
		location.x -= 1

	case DirectionRight:
		location.x += 1
	}

	width, height := termbox.Size()
	if location.x <= 0 || location.x >= width ||
		location.y <= 0 || location.y >= height {
		bullet.alive = false
	}

	return nil
}
