// common
// @author xiangqian
// @date 10:58 2022/11/27
package tank

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten"
	"log"
)

// type
type (
	Direction int8  // 方向
	Speed     int16 // 速度
	Status    int8  // 状态
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

// 状态
const (
	StatusNew  Status = iota // 初始化
	StatusRun                // 运行（执行）
	StatusTerm               // 终止
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// 位置信息
type Location struct {
	x float64
	y float64
}

// 图形
type Graphics interface {
	// id
	Id() string

	// 方向
	Direction() Direction

	// 速度
	Speed() Speed

	// 状态
	Status() Status

	// 生命值
	HP() int8

	// 绘制
	Draw(screen *ebiten.Image) error
}

// 图形
type AbsGraphics struct {
	id        string    // id
	direction Direction // 方向
	speed     Speed     // 速度
	status    Status    // 状态
	hp        int8      // 生命值
}

func (absGraphics *AbsGraphics) Id() string {
	return absGraphics.id
}

func (absGraphics *AbsGraphics) Direction() Direction {
	return absGraphics.direction
}

func (absGraphics *AbsGraphics) SetDirection(direction Direction) {
	absGraphics.direction = direction
}

func (absGraphics *AbsGraphics) Speed() Speed {
	return absGraphics.speed
}

func (absGraphics *AbsGraphics) Status() Status {
	return absGraphics.status
}

func (absGraphics *AbsGraphics) HP() int8 {
	return absGraphics.hp
}

func (absGraphics *AbsGraphics) Draw(screen *ebiten.Image) error {
	if absGraphics.status == StatusTerm {
		return errors.New(fmt.Sprintf("the %v has been terminated", absGraphics.Id()))
	}
	return nil
}

// 越界校验
func (absGraphics *AbsGraphics) OutOfBoundsCheck() bool {
	return false
}

func Uuid() string {
	return uuid.New().String()
}
