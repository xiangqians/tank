// Graphics
//
// @author xiangqian
// @date 00:20 2022/12/02
package tank

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"reflect"
)

// type
type (
	Direction  int8  // 方向
	Speed      int16 // 速度
	Status     int8  // 状态
	GraphicsTy int8  // 图像类型
)

// 方向
const (
	DirectionUp Direction = iota + 1
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
	StatusNew  Status = iota + 1 // 初始化
	StatusRun                    // 运行（执行）
	StatusTerm                   // 终止
)

const (
	GraphicsTyTank GraphicsTy = iota + 1
	GraphicsTyBullet
)

// 位置信息
type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// 图形
type Graphics interface {
	// id
	GetId() string

	// 状态
	GetStatus() Status

	// 生命值
	GetHp() uint8

	// 绘制
	Draw(screen *ebiten.Image) error
}

// 图形
type AbsGraphics struct {
	Id         string        `json:"id"`         // id
	GraphicsTy GraphicsTy    `json:"graphicsTy"` // 图像类型
	Location   Location      `json:"location"`   // 位置
	Direction  Direction     `json:"direction"`  // 方向
	Speed      Speed         `json:"speed"`      // 速度
	Status     Status        `json:"status"`     // 状态
	Hp         uint8         `json:"hp"`         // 生命值
	pImg       *ebiten.Image // 图片
	sub        interface{}   // 子类
}

func CreateAbsGraphics(graphicsTy GraphicsTy, location Location, direction Direction, speed Speed) *AbsGraphics {
	return &AbsGraphics{
		Id:         Uuid(),
		GraphicsTy: graphicsTy,
		Location:   location,
		Direction:  direction,
		Speed:      speed,
		Status:     StatusNew,
		Hp:         100,
		pImg:       nil,
		sub:        nil,
	}
}

func (pAbsGraphics *AbsGraphics) Init(sub interface{}) {
	pAbsGraphics.sub = sub
	pAbsGraphics.pImg = pAbsGraphics.DirectionImg(pAbsGraphics.Direction)
}

func (pAbsGraphics *AbsGraphics) GetId() string {
	return pAbsGraphics.Id
}

func (pAbsGraphics *AbsGraphics) GetStatus() Status {
	return pAbsGraphics.Status
}

func (pAbsGraphics *AbsGraphics) GetHp() uint8 {
	return pAbsGraphics.Hp
}

func (pAbsGraphics *AbsGraphics) Draw(screen *ebiten.Image) error {
	if pAbsGraphics.Status == StatusTerm {
		return errors.New(fmt.Sprintf("the %v has been terminated", pAbsGraphics.Id))
	}

	op := &ebiten.DrawImageOptions{}
	location := pAbsGraphics.Location
	op.GeoM.Translate(location.X, location.Y)
	return screen.DrawImage(pAbsGraphics.pImg, op)
}

func (pAbsGraphics *AbsGraphics) Move(direction Direction) {
	// -→ x
	// ↓ y
	pLocation := &pAbsGraphics.Location
	pAbsGraphics.Direction = direction
	var pImg *ebiten.Image
	var xx float64 = 1 + float64(pAbsGraphics.Speed)
	switch direction {
	case DirectionUp:
		newy := pLocation.Y - xx
		if !pAbsGraphics.refIsOutOfBounds(pLocation.X, newy) {
			pLocation.Y = newy
		}
		pImg = pAbsGraphics.refUpImg()

	case DirectionDown:
		newy := pLocation.Y + xx
		if !pAbsGraphics.refIsOutOfBounds(pLocation.X, newy) {
			pLocation.Y = newy
		}
		pImg = pAbsGraphics.refDownImg()

	case DirectionLeft:
		newx := pLocation.X - xx
		if !pAbsGraphics.refIsOutOfBounds(newx, pLocation.Y) {
			pLocation.X = newx
		}
		pImg = pAbsGraphics.refLeftImg()

	case DirectionRight:
		newx := pLocation.X + xx
		if !pAbsGraphics.refIsOutOfBounds(newx, pLocation.Y) {
			pLocation.X = newx
		}
		pImg = pAbsGraphics.refRightImg()
	}

	pAbsGraphics.pImg = pImg
}

func (pAbsGraphics *AbsGraphics) refIsOutOfBounds(x, y float64) bool {
	r := pAbsGraphics.refMethod("IsOutOfBounds", []reflect.Value{reflect.ValueOf(x), reflect.ValueOf(y)})
	if r != nil {
		return r.(bool)
	}

	//panic(nil)
	return pAbsGraphics.IsOutOfBounds(x, y)
}

// 是否越界
func (pAbsGraphics *AbsGraphics) IsOutOfBounds(x, y float64) bool {
	width, height := pAbsGraphics.pImg.Size()
	if x <= 0 || x >= screenWidth-float64(height) || y <= 0 || y >= screenHeight-float64(width) {
		return true
	}

	return false
}

func (pAbsGraphics *AbsGraphics) DirectionImg(direction Direction) *ebiten.Image {
	var pImg *ebiten.Image = nil
	switch direction {
	case DirectionUp:
		pImg = pAbsGraphics.refUpImg()

	case DirectionDown:
		pImg = pAbsGraphics.refDownImg()

	case DirectionLeft:
		pImg = pAbsGraphics.refLeftImg()

	case DirectionRight:
		pImg = pAbsGraphics.refRightImg()
	}
	return pImg
}

func (pAbsGraphics *AbsGraphics) refUpImg() *ebiten.Image {
	return pAbsGraphics.refNameImg("UpImg")
}

func (pAbsGraphics *AbsGraphics) refDownImg() *ebiten.Image {
	return pAbsGraphics.refNameImg("DownImg")
}

func (pAbsGraphics *AbsGraphics) refLeftImg() *ebiten.Image {
	return pAbsGraphics.refNameImg("LeftImg")
}

func (pAbsGraphics *AbsGraphics) refRightImg() *ebiten.Image {
	return pAbsGraphics.refNameImg("RightImg")
}

func (pAbsGraphics *AbsGraphics) UpImg() *ebiten.Image {
	panic(nil)
}

func (pAbsGraphics *AbsGraphics) DownImg() *ebiten.Image {
	panic(nil)
}

func (pAbsGraphics *AbsGraphics) LeftImg() *ebiten.Image {
	panic(nil)
}

func (pAbsGraphics *AbsGraphics) RightImg() *ebiten.Image {
	panic(nil)
}

func (pAbsGraphics *AbsGraphics) refNameImg(name string) *ebiten.Image {
	r := pAbsGraphics.refMethod(name, nil)
	if r != nil {
		return r.(*ebiten.Image)
	}

	panic(nil)
}

// 反射执行方法
// name: 方法名
// in: 入参，如果没有参数可以传 nil 或者空切片 make([]reflect.Value, 0)
func (pAbsGraphics *AbsGraphics) refMethod(name string, in []reflect.Value) interface{} {
	ref := reflect.ValueOf(pAbsGraphics.sub)
	method := ref.MethodByName(name)
	if method.IsValid() {
		r := method.Call(in)
		return r[0].Interface()
	}
	return nil
}
