// tank
// 坦克
//
// @author xiangqian
// @date 10:58 2022/11/27
package tank

import (
	"github.com/nsf/termbox-go"
)

type Tank struct {
	*AbsGraphics

	body []*Location // 坦克体
}

func CreateTank(location Location, direction Direction, speed Speed) *Tank {
	pTank := &Tank{
		AbsGraphics: &AbsGraphics{
			id:        "tank_" + Uuid(),
			direction: direction,
			speed:     speed,
			status:    StatusNew,
		},
		body: CreateTankBody(location, direction),
	}
	return pTank
}

// location 中心位置
func CreateTankBody(location Location, direction Direction) []*Location {
	x := location.x
	y := location.y
	switch direction {
	case DirectionUp:
		y -= 1
		return []*Location{
			// 头部
			{x, y},
			// 左侧
			{x - 2, y + 1},
			{x - 2, y + 2},
			// 右侧
			{x + 2, y + 1},
			{x + 2, y + 2},
		}

	case DirectionDown:
		y += 1
		return []*Location{
			// 头部
			{x, y},
			// 左侧
			{x - 2, y - 2},
			{x - 2, y - 1},
			// 右侧
			{x + 2, y - 2},
			{x + 2, y - 1},
		}

	case DirectionLeft:
		x -= 2
		return []*Location{
			// 头部
			{x, y},
			// 上侧
			{x + 2, y - 1},
			{x + 4, y - 1},
			// 下侧
			{x + 2, y + 1},
			{x + 4, y + 1},
		}

	case DirectionRight:
		x += 2
		return []*Location{
			// 头部
			{x, y},
			// 上侧
			{x - 2, y - 1},
			{x - 4, y - 1},
			// 下侧
			{x - 2, y + 1},
			{x - 4, y + 1},
		}
	}

	return nil
}

// 中心位置 location
func (tank *Tank) Location() Location {
	location := *tank.body[0]
	x, y := location.x, location.y
	switch tank.direction {
	case DirectionUp:
		y += 1

	case DirectionDown:
		y -= 1

	case DirectionLeft:
		x += 2

	case DirectionRight:
		x -= 2
	}

	return Location{x, y}
}

func (tank *Tank) Move(direction Direction) {

	if tank.direction != direction {
		tank.body = CreateTankBody(tank.Location(), direction)
	}
	tank.direction = direction

	location := *tank.body[0]
	width, height := termbox.Size()
	if location.x <= infoBar.width+2 || location.x+2 >= width ||
		location.y <= 0 || location.y+1 >= height {
		return
	}

	body := tank.body
	switch direction {
	case DirectionUp:
		for _, location := range body {
			location.y -= 1
		}

	case DirectionDown:
		for _, location := range body {
			location.y += 1
		}

	case DirectionLeft:
		for _, location := range body {
			location.x -= 1
		}

	case DirectionRight:
		for _, location := range body {
			location.x += 1
		}
	}

}

// 开火
func (tank *Tank) Fire() {
	addGraphics(CreateBullet(tank, SpeedNormal))
}

func (tank *Tank) Draw() error {
	if err := tank.AbsGraphics.Draw(); err != nil {
		return err
	}

	body := tank.body
	for _, location := range body {
		termbox.SetCell(location.x, location.y, '⬛', termbox.ColorRed, termbox.ColorRed)
	}

	return nil
}
