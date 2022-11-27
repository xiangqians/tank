// common
// @author xiangqian
// @date 10:58 2022/11/27
package tank

import (
	"errors"
	"github.com/google/uuid"
)

// type
type (
	Direction int8 // 方向
	Speed     int8 // 速度
)

// 方向
const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

// 速度
const (
	SpeedSlow Speed = iota
	SpeedNormal
	SpeedFast
)

// 位置信息
type Location struct {
	x int
	y int
}

// 图形
type Graphics interface {
	// 是否存活
	Alive() bool

	// 绘制
	Draw() error
}

// 图形
type AbsGraphics struct {
	id        string    // id
	direction Direction // 方向
	speed     Speed     // 速度
	alive     bool      // 是否存活
}

func (absGraphics *AbsGraphics) Direction() Direction {
	return absGraphics.direction
}

func (absGraphics *AbsGraphics) Alive() bool {
	return absGraphics.alive
}

func (absGraphics *AbsGraphics) Draw() error {
	if !absGraphics.alive {
		return errors.New("the bullet not alive")
	}
	return nil
}

func Uuid() string {
	return uuid.New().String()
}
