// Graphics
//
// @author xiangqian
// @date 00:20 2022/12/02
package tank

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"log"
	"reflect"
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
	SpeedSlow Speed = iota * 2
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

	// 状态
	Status() Status

	// 生命值
	HP() uint8

	// 绘制
	Draw(screen *ebiten.Image) error
}

// 图形
type AbsGraphics struct {
	id        string        // id
	pLocation *Location     // 位置
	pImg      *ebiten.Image // 图片
	direction Direction     // 方向
	speed     Speed         // 速度
	status    Status        // 状态
	hp        uint8         // 生命值
	sub       interface{}   // 子类
}

func CreateAbsGraphics(location Location, direction Direction, speed Speed) *AbsGraphics {
	return &AbsGraphics{
		id:        Uuid(),
		pLocation: &Location{location.x, location.y},
		pImg:      pDefaultImg,
		direction: direction,
		speed:     speed,
		status:    StatusNew,
		hp:        100,
	}
}

func (absGraphics *AbsGraphics) Init(sub interface{}) {
	absGraphics.sub = sub
	absGraphics.pImg = absGraphics.DirectionImg(absGraphics.direction)
}

func (absGraphics *AbsGraphics) Id() string {
	return absGraphics.id
}

func (absGraphics *AbsGraphics) Status() Status {
	return absGraphics.status
}

func (absGraphics *AbsGraphics) HP() uint8 {
	return absGraphics.hp
}

func (absGraphics *AbsGraphics) Draw(screen *ebiten.Image) error {
	if absGraphics.status == StatusTerm {
		return errors.New(fmt.Sprintf("the %v has been terminated", absGraphics.Id()))
	}

	op := &ebiten.DrawImageOptions{}
	location := absGraphics.pLocation
	op.GeoM.Translate(location.x, location.y)
	return screen.DrawImage(absGraphics.pImg, op)
}

func (absGraphics *AbsGraphics) Move(direction Direction) {
	// -→ x
	// ↓ y
	location := absGraphics.pLocation
	absGraphics.direction = direction
	var pImg *ebiten.Image
	var xx float64 = 1 + float64(absGraphics.speed)
	switch direction {
	case DirectionUp:
		newy := location.y - xx
		if !absGraphics.refIsOutOfBounds(location.x, newy) {
			location.y = newy
		}
		pImg = absGraphics.refUpImg()

	case DirectionDown:
		newy := location.y + xx
		if !absGraphics.refIsOutOfBounds(location.x, newy) {
			location.y = newy
		}
		pImg = absGraphics.refDownImg()

	case DirectionLeft:
		newx := location.x - xx
		if !absGraphics.refIsOutOfBounds(newx, location.y) {
			location.x = newx
		}
		pImg = absGraphics.refLeftImg()

	case DirectionRight:
		newx := location.x + xx
		if !absGraphics.refIsOutOfBounds(newx, location.y) {
			location.x = newx
		}
		pImg = absGraphics.refRightImg()
	}

	absGraphics.pImg = pImg
}

func (absGraphics *AbsGraphics) refIsOutOfBounds(x, y float64) bool {
	r := absGraphics.refMethod("IsOutOfBounds", []reflect.Value{reflect.ValueOf(x), reflect.ValueOf(y)})
	if r != nil {
		return r.(bool)
	}

	//panic(nil)
	return absGraphics.IsOutOfBounds(x, y)
}

// 是否越界
func (absGraphics *AbsGraphics) IsOutOfBounds(x, y float64) bool {
	width, height := absGraphics.pImg.Size()
	if x <= 0 || x >= screenWidth-float64(height) || y <= 0 || y >= screenHeight-float64(width) {
		return true
	}

	return false
}

func (absGraphics *AbsGraphics) DirectionImg(direction Direction) *ebiten.Image {
	var pImg *ebiten.Image = nil
	switch direction {
	case DirectionUp:
		pImg = absGraphics.refUpImg()

	case DirectionDown:
		pImg = absGraphics.refDownImg()

	case DirectionLeft:
		pImg = absGraphics.refLeftImg()

	case DirectionRight:
		pImg = absGraphics.refRightImg()
	}
	return pImg
}

func (absGraphics *AbsGraphics) refUpImg() *ebiten.Image {
	return absGraphics.refNameImg("UpImg")
}

func (absGraphics *AbsGraphics) refDownImg() *ebiten.Image {
	return absGraphics.refNameImg("DownImg")
}

func (absGraphics *AbsGraphics) refLeftImg() *ebiten.Image {
	return absGraphics.refNameImg("LeftImg")
}

func (absGraphics *AbsGraphics) refRightImg() *ebiten.Image {
	return absGraphics.refNameImg("RightImg")
}

func (absGraphics *AbsGraphics) UpImg() *ebiten.Image {
	panic(nil)
}

func (absGraphics *AbsGraphics) DownImg() *ebiten.Image {
	panic(nil)
}

func (absGraphics *AbsGraphics) LeftImg() *ebiten.Image {
	panic(nil)
}

func (absGraphics *AbsGraphics) RightImg() *ebiten.Image {
	panic(nil)
}

func (absGraphics *AbsGraphics) refNameImg(name string) *ebiten.Image {
	r := absGraphics.refMethod(name, nil)
	if r != nil {
		return r.(*ebiten.Image)
	}

	panic(nil)
}

// 反射执行方法
// name: 方法名
// in: 入参，如果没有参数可以传 nil 或者空切片 make([]reflect.Value, 0)
func (absGraphics *AbsGraphics) refMethod(name string, in []reflect.Value) interface{} {
	ref := reflect.ValueOf(absGraphics.sub)
	method := ref.MethodByName(name)
	if method.IsValid() {
		r := method.Call(in)
		return r[0].Interface()
	}
	return nil
}
