// bullet
// 子弹
//
// @author xiangqian
// @date 11:01 2022/11/27
package tank

import (
	"github.com/nsf/termbox-go"
	"time"
)

type Bullet struct {
	*AbsGraphics
	tankId   string    // 子弹所属坦克id
	location *Location // 子弹位置
}

func CreateBullet(tank *Tank, speed Speed) *Bullet {
	location := tank.Location()
	return &Bullet{
		AbsGraphics: &AbsGraphics{
			id:        "bullet_" + Uuid(),
			direction: tank.direction,
			speed:     speed,
			status:    StatusNew,
		},
		tankId:   tank.id,
		location: &Location{location.x, location.y},
	}
}
func (bullet *Bullet) check() {
	width, height := termbox.Size()
	location := bullet.location
	if location.x <= 0 || location.x >= width ||
		location.y <= 0 || location.y >= height {
		bullet.status = StatusTerm
	}
}

func (bullet *Bullet) Draw() error {
	if err := bullet.AbsGraphics.Draw(); err != nil {
		return err
	}

	location := bullet.location
	termbox.SetCell(location.x, location.y, '⬛', termbox.ColorRed, termbox.ColorRed)

	// 如果是当前用户所发射的子弹，那么由当前用户设置子弹
	if bullet.tankId == tank.id && bullet.status == StatusNew {
		go func() {
			bullet.status = StatusRun
			for {
				bullet.check()
				if bullet.status != StatusRun {
					break
				}
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

				switch bullet.speed {
				case SpeedSlow:
					time.Sleep(100 * time.Millisecond)

				case SpeedNormal:
					time.Sleep(50 * time.Millisecond)

				case SpeedFast:
					time.Sleep(10 * time.Millisecond)
				}
			}
			//fmt.Printf("%v Tserm\n", bullet.id)
		}()
	}

	bullet.check()

	return nil
}
