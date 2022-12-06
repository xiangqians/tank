// Graphics
//
// @author xiangqian
// @date 00:20 2022/12/02
package tank

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"math"
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

	GetAbsGraphics() *AbsGraphics
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
	pImage     *ebiten.Image // 图片
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
		pImage:     nil,
		sub:        nil,
	}
}

func (pAbsGraphics *AbsGraphics) Init(sub interface{}) {
	pAbsGraphics.sub = sub
	pAbsGraphics.pImage = pAbsGraphics.DirectionImage(pAbsGraphics.Direction)
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
	return screen.DrawImage(pAbsGraphics.pImage, op)
}

func (pAbsGraphics *AbsGraphics) GetAbsGraphics() *AbsGraphics {
	return pAbsGraphics
}

func (pAbsGraphics *AbsGraphics) Move(direction Direction) {
	// -→ x
	// ↓ y
	pLocation := &pAbsGraphics.Location
	pAbsGraphics.Direction = direction
	var pImage *ebiten.Image
	var xx float64 = 1 + float64(pAbsGraphics.Speed)
	switch direction {
	case DirectionUp:
		newy := pLocation.Y - xx
		if !pAbsGraphics.callIsOutOfBounds(pLocation.X, newy) {
			pLocation.Y = newy
		}
		pImage = pAbsGraphics.callUpImage()

	case DirectionDown:
		newy := pLocation.Y + xx
		if !pAbsGraphics.callIsOutOfBounds(pLocation.X, newy) {
			pLocation.Y = newy
		}
		pImage = pAbsGraphics.callDownImage()

	case DirectionLeft:
		newx := pLocation.X - xx
		if !pAbsGraphics.callIsOutOfBounds(newx, pLocation.Y) {
			pLocation.X = newx
		}
		pImage = pAbsGraphics.callLeftImage()

	case DirectionRight:
		newx := pLocation.X + xx
		if !pAbsGraphics.callIsOutOfBounds(newx, pLocation.Y) {
			pLocation.X = newx
		}
		pImage = pAbsGraphics.callRightImage()
	}

	pAbsGraphics.pImage = pImage
}

// 判断图形是否相交
func (pAbsGraphics *AbsGraphics) Intersect(x, y float64, pOtherAbsGraphics *AbsGraphics) bool {
	// 两个矩形相交机几种情况：images/rectangle_itersect.png
	// 重心距离在X轴和Y轴都小于两矩形的长或宽的一半之和

	width, height := pAbsGraphics.pImage.Size()
	//centerX := pAbsGraphics.Location.X + float64(width/2)
	//centerY := pAbsGraphics.Location.Y + float64(height/2)
	centerX := x + float64(width/2)
	centerY := y + float64(height/2)
	//log.Printf("center x: %v, y: %v\n", centerX, centerY)

	otherWidth, otherHeight := pOtherAbsGraphics.pImage.Size()
	otherCenterX := pOtherAbsGraphics.Location.X + float64(otherWidth/2)
	otherCenterY := pOtherAbsGraphics.Location.Y + float64(otherHeight/2)
	//log.Printf("otherCenter x: %v, y: %v\n", otherCenterX, otherCenterY)

	centerWidth := math.Abs(centerX - otherCenterX)
	centerHeight := math.Abs(centerY - otherCenterY)
	//log.Printf("center width: %v, height: %v\n", centerWidth, centerHeight)

	if centerWidth <= float64((width+otherWidth)/2) && centerHeight <= float64((height+otherHeight)/2) {
		return true
	}

	return false
}

func (pAbsGraphics *AbsGraphics) callIsOutOfBounds(x, y float64) bool {
	r := CallMethod(pAbsGraphics.sub, "IsOutOfBounds", []reflect.Value{reflect.ValueOf(x), reflect.ValueOf(y)})
	if r != nil {
		return r.(bool)
	}

	//panic(nil)
	return pAbsGraphics.IsOutOfBounds(x, y)
}

// 是否越界
func (pAbsGraphics *AbsGraphics) IsOutOfBounds(x, y float64) bool {
	width, height := pAbsGraphics.pImage.Size()
	if x <= 0 || x >= screenWidth-float64(height) || y <= 0 || y >= screenHeight-float64(width) {
		return true
	}

	for _, value := range pApp.pGame.GraphicsMap {
		if value.GetId() == pAbsGraphics.GetId() {
			continue
		}
		if pAbsGraphics.Intersect(x, y, value.GetAbsGraphics()) {
			return true
		}
	}

	return false
}

func (pAbsGraphics *AbsGraphics) DirectionImage(direction Direction) *ebiten.Image {
	var pImage *ebiten.Image = nil
	switch direction {
	case DirectionUp:
		pImage = pAbsGraphics.callUpImage()

	case DirectionDown:
		pImage = pAbsGraphics.callDownImage()

	case DirectionLeft:
		pImage = pAbsGraphics.callLeftImage()

	case DirectionRight:
		pImage = pAbsGraphics.callRightImage()
	}
	return pImage
}

func (pAbsGraphics *AbsGraphics) callUpImage() *ebiten.Image {
	return pAbsGraphics.callNameImage("UpImage")
}

func (pAbsGraphics *AbsGraphics) callDownImage() *ebiten.Image {
	return pAbsGraphics.callNameImage("DownImage")
}

func (pAbsGraphics *AbsGraphics) callLeftImage() *ebiten.Image {
	return pAbsGraphics.callNameImage("LeftImage")
}

func (pAbsGraphics *AbsGraphics) callRightImage() *ebiten.Image {
	return pAbsGraphics.callNameImage("RightImage")
}

func (pAbsGraphics *AbsGraphics) UpImage() *ebiten.Image {
	panic(nil)
}

func (pAbsGraphics *AbsGraphics) DownImage() *ebiten.Image {
	panic(nil)
}

func (pAbsGraphics *AbsGraphics) LeftImage() *ebiten.Image {
	panic(nil)
}

func (pAbsGraphics *AbsGraphics) RightImage() *ebiten.Image {
	panic(nil)
}

func (pAbsGraphics *AbsGraphics) callNameImage(name string) *ebiten.Image {
	r := CallMethod(pAbsGraphics.sub, name, nil)
	if r != nil {
		return r.(*ebiten.Image)
	}

	panic(nil)
}
