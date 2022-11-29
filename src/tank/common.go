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
	id        string                   // id
	pLocation *Location                // 位置
	pOptions  *ebiten.DrawImageOptions // 绘制图片操作类
	pImg      *ebiten.Image            // 图片
	direction Direction                // 方向
	speed     Speed                    // 速度
	status    Status                   // 状态
	hp        uint8                    // 生命值
	sub       interface{}              // 子类
}

func CreateAbsGraphics(location Location, direction Direction, speed Speed) *AbsGraphics {
	return &AbsGraphics{
		id:        Uuid(),
		pLocation: &Location{location.x, location.y},
		pOptions:  nil,
		pImg:      nil,
		direction: direction,
		speed:     speed,
		status:    StatusNew,
		hp:        100,
	}
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

	options := absGraphics.pOptions
	if options == nil {
		options = &ebiten.DrawImageOptions{}
		absGraphics.pOptions = options
	}

	location := absGraphics.pLocation
	options.GeoM.Reset()
	options.GeoM.Translate(location.x, location.y)
	pImg := absGraphics.pImg
	if pImg == nil {
		pImg = absGraphics.DirectionImg(absGraphics.direction)
		absGraphics.pImg = pImg
	}
	return screen.DrawImage(pImg, options)
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
		if !absGraphics.IsOutOfBounds(location.x, newy) {
			location.y = newy
			pImg = absGraphics.UpImg()
		}

	case DirectionDown:
		newy := location.y + xx
		if !absGraphics.IsOutOfBounds(location.x, newy) {
			location.y = newy
			pImg = absGraphics.DownImg()
		}

	case DirectionLeft:
		newx := location.x - xx
		if !absGraphics.IsOutOfBounds(newx, location.y) {
			location.x = newx
			pImg = absGraphics.LeftImg()
		}

	case DirectionRight:
		newx := location.x + xx
		if !absGraphics.IsOutOfBounds(newx, location.y) {
			location.x = newx
			pImg = absGraphics.RightImg()
		}
	}
	absGraphics.pImg = pImg
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
		pImg = absGraphics.UpImg()

	case DirectionDown:
		pImg = absGraphics.DownImg()

	case DirectionLeft:
		pImg = absGraphics.LeftImg()

	case DirectionRight:
		pImg = absGraphics.RightImg()
	}
	return pImg
}

func (absGraphics *AbsGraphics) UpImg() *ebiten.Image {
	return absGraphics.NameImg("UpImg")
}

func (absGraphics *AbsGraphics) DownImg() *ebiten.Image {
	return absGraphics.NameImg("DownImg")
}

func (absGraphics *AbsGraphics) LeftImg() *ebiten.Image {
	return absGraphics.NameImg("LeftImg")
}

func (absGraphics *AbsGraphics) RightImg() *ebiten.Image {
	return absGraphics.NameImg("RightImg")
}

func (absGraphics *AbsGraphics) NameImg(name string) *ebiten.Image {
	ref := reflect.ValueOf(absGraphics.sub)
	method := ref.MethodByName(name)
	if method.IsValid() {
		r := method.Call(make([]reflect.Value, 0))
		return r[0].Interface().(*ebiten.Image)
	}
	panic(nil)
}

func Uuid() string {
	return uuid.New().String()
}
