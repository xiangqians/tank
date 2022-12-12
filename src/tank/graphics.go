// Graphics
//
// @author xiangqian
// @date 00:20 2022/12/02
package tank

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"math"
	"reflect"
	"time"
)

// type
type (
	Direction  int8 // 方向
	Speed      int8 // 速度
	Status     int8 // 状态
	GraphicsTy int8 // 图像类型
)

const DefaultHp uint8 = 32

// 方向
const (
	DirectionUp Direction = iota + 1
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
	StatusNew  Status = iota + 1 // 初始化
	StatusRun                    // 运行（执行）
	StatusTerm                   // 终止
)

const (
	GraphicsTyTank GraphicsTy = iota + 1
	GraphicsTyBullet
	GraphicsTyEquip
)

// 位置信息
type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// 坦克元数据字体
// Metadata
var tankMdFont font.Face

func init() {
	tankMdFont = CreateFontFace(8, 72)
}

func GetSpeedByName(name string) (Speed, error) {
	switch name {
	case "Slow":
		return SpeedSlow, nil
	case "Normal":
		return SpeedNormal, nil
	case "Fast":
		return SpeedFast, nil
	}
	return 0, errors.New("unknown name")
}

func SpeedString(speed Speed) (string, error) {
	switch speed {
	case SpeedSlow:
		return "Slow", nil
	case SpeedNormal:
		return "Normal", nil
	case SpeedFast:
		return "Fast", nil
	}
	return "", errors.New("unknown value")
}

// 图形
type Graphics interface {
	// id
	GetId() string

	// name
	GetName() string

	// 获取图像类型
	GetGraphicsTy() GraphicsTy

	// 获取位置
	GetLocation() Location

	// 状态
	GetStatus() Status

	// 生命值
	GetHp() uint8

	// 时间戳
	GetTimestamp() int64

	// 校验时间戳是否有效
	VerifyTimestamp() bool

	// 获取图片
	GetImage() *ebiten.Image

	// 绘制
	Draw(screen *ebiten.Image) error
}

func DeserializeBytesToGraphics(bytes []byte) Graphics {
	pAbsGraphics := &AbsGraphics{}
	err := Deserialize(bytes, pAbsGraphics)
	if err != nil {
		return nil
	}

	var graphics Graphics = nil
	switch pAbsGraphics.GraphicsTy {
	case GraphicsTyTank:
		//pTank := &Tank{AbsGraphics: pAbsGraphics}
		//pTank.Timestamp = time.Now().UnixNano()
		//pTank.Init()
		//graphics = pTank

		pTank := &Tank{}
		err := Deserialize(bytes, pTank)
		if err == nil {
			pTank.Timestamp = time.Now().UnixNano()
			pTank.Init()
			graphics = pTank
		}

	case GraphicsTyBullet:
		pBullet := &Bullet{AbsGraphics: pAbsGraphics}
		pBullet.Timestamp = time.Now().UnixNano()
		pBullet.Init()
		graphics = pBullet

	case GraphicsTyEquip:
		pEquip := &Equip{}
		err := Deserialize(bytes, pEquip)
		if err == nil {
			pEquip.Timestamp = time.Now().UnixNano()
			pEquip.Init()
			graphics = pEquip
		}
	}

	return graphics
}

// 图形
type AbsGraphics struct {
	Id         string        `json:"id"`         // id
	Name       string        `json:"name"`       // 名称
	GraphicsTy GraphicsTy    `json:"graphicsTy"` // 图像类型
	Location   Location      `json:"location"`   // 位置
	Direction  Direction     `json:"direction"`  // 方向
	Speed      Speed         `json:"speed"`      // 速度
	Status     Status        `json:"status"`     // 状态
	Hp         uint8         `json:"hp"`         // 生命值
	Timestamp  int64         `json:"timestamp"`  // 时间戳，ns
	pImage     *ebiten.Image // 图片
	sub        interface{}   // 子类
}

func CreateAbsGraphics(name string, graphicsTy GraphicsTy, location Location, direction Direction, speed Speed) *AbsGraphics {
	return &AbsGraphics{
		Id:         Uuid(),
		Name:       name,
		GraphicsTy: graphicsTy,
		Location:   location,
		Direction:  direction,
		Speed:      speed,
		Status:     StatusNew,
		Hp:         DefaultHp,
		Timestamp:  time.Now().UnixNano(),
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

func (pAbsGraphics *AbsGraphics) GetName() string {
	return pAbsGraphics.Name
}

func (pAbsGraphics *AbsGraphics) GetGraphicsTy() GraphicsTy {
	return pAbsGraphics.GraphicsTy
}

func (pAbsGraphics *AbsGraphics) GetLocation() Location {
	return pAbsGraphics.Location
}

func (pAbsGraphics *AbsGraphics) GetStatus() Status {
	return pAbsGraphics.Status
}

func (pAbsGraphics *AbsGraphics) GetHp() uint8 {
	return pAbsGraphics.Hp
}

func (pAbsGraphics *AbsGraphics) GetTimestamp() int64 {
	return pAbsGraphics.Timestamp
}

func (pAbsGraphics *AbsGraphics) VerifyTimestamp() bool {
	// 不对终止图形校验
	if pAbsGraphics.Status == StatusTerm {
		return false
	}

	// 只校验子弹
	if pAbsGraphics.GraphicsTy != GraphicsTyBullet {
		return true
	}

	pBullet := reflect.ValueOf(pAbsGraphics.sub).Interface().(*Bullet)

	// 不校验属于当前坦克所发射出的子弹
	if pBullet.TankId == pApp.pGame.pTank.Id {
		return true
	}

	// 不属于当前坦克所发射出的子弹超时（2s）未动则视为终止
	if time.Now().UnixNano()-pBullet.Timestamp > 2*int64(time.Second) {
		pBullet.Status = StatusTerm
		return false
	}

	return true
}

func (pAbsGraphics *AbsGraphics) GetImage() *ebiten.Image {
	return pAbsGraphics.pImage
}

func (pAbsGraphics *AbsGraphics) Draw(screen *ebiten.Image) error {
	if pAbsGraphics.Status == StatusTerm {
		return errors.New(fmt.Sprintf("the %v has been terminated", pAbsGraphics.Id))
	}

	// 不绘制隐形坦克（除了当前坦克外）
	if pAbsGraphics.GraphicsTy == GraphicsTyTank &&
		pAbsGraphics.Id != pApp.pGame.pTank.Id &&
		reflect.ValueOf(pAbsGraphics.sub).Interface().(*Tank).TankInvisFlag {
		return nil
	}

	// 绘制图形
	op := &ebiten.DrawImageOptions{}
	location := pAbsGraphics.Location
	op.GeoM.Translate(location.X, location.Y)
	err := screen.DrawImage(pAbsGraphics.pImage, op)

	// 绘制装备类型
	if pAbsGraphics.GraphicsTy == GraphicsTyEquip {
		x, y := int(location.X), int(location.Y)
		y -= 6

		t := ""
		switch pAbsGraphics.sub.(*Equip).EquipType {
		// 坦克加速
		case EquipTypeTankAcc:
			t = "TS"

		// 子弹加速
		case EquipTypeBulletAcc:
			t = "BS"

		// HP增加
		case EquipTypeHpInc:
			t = "HI"

		// 坦克隐形
		case EquipTypeTankInvis:
			t = "TI"

		default:
			t = "UK"
		}
		text.Draw(screen, t, tankMdFont, x, y, color.White)
	}

	// 绘制坦克元数据（除了当前坦克外）
	if err == nil &&
		pAbsGraphics.Id != pApp.pGame.pTank.Id &&
		pAbsGraphics.GraphicsTy == GraphicsTyTank {
		nameX, hpX := int(location.X), int(location.X)
		nameY, hpY := int(location.Y), int(location.Y)
		switch pAbsGraphics.Direction {
		case DirectionUp:
			_, height := pAbsGraphics.pImage.Size()
			nameY += height + 10
			hpY = nameY
			hpY += 10
		case DirectionDown, DirectionLeft, DirectionRight:
			hpY -= 8
			nameY = hpY
			nameY -= 10
		}
		text.Draw(screen, fmt.Sprintf("%v", pAbsGraphics.GetName()), tankMdFont, nameX, nameY, color.White)
		text.Draw(screen, fmt.Sprintf("HP: %v", pAbsGraphics.GetHp()), tankMdFont, hpX, hpY, color.White)
	}

	return err
}

func (pAbsGraphics *AbsGraphics) GetAbsGraphics() *AbsGraphics {
	return pAbsGraphics
}

func (pAbsGraphics *AbsGraphics) Move(direction Direction) {

	// Speed value
	var speedValue float64 = 1
	switch pAbsGraphics.Speed {
	case SpeedSlow:
		speedValue += 1
	case SpeedNormal:
		speedValue += 3
	case SpeedFast:
		speedValue += 5
	}

	// -→ x
	// ↓ y
	pLocation := &pAbsGraphics.Location
	pAbsGraphics.Direction = direction
	var pImage *ebiten.Image
	var xx float64 = 1 + speedValue
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

func (pAbsGraphics *AbsGraphics) callIntersect(x, y float64, otherGraphics Graphics) bool {
	r := CallMethod(pAbsGraphics.sub, "Intersect", []reflect.Value{reflect.ValueOf(x), reflect.ValueOf(y), reflect.ValueOf(otherGraphics)})
	if r != nil {
		return r.(bool)
	}

	return pAbsGraphics.Intersect(x, y, otherGraphics)
}

// 判断图形是否相交
func (pAbsGraphics *AbsGraphics) Intersect(x, y float64, otherGraphics Graphics) bool {

	// 当前主动校验的是tank时，不校验tank自身发射的子弹
	if pAbsGraphics.GraphicsTy == GraphicsTyTank &&
		otherGraphics.GetGraphicsTy() == GraphicsTyBullet &&
		otherGraphics.(*Bullet).TankId == pAbsGraphics.Id {
		return false
	}

	// 当前主动校验的是bullet时，不对tank有效
	if pAbsGraphics.GraphicsTy == GraphicsTyBullet &&
		otherGraphics.GetGraphicsTy() == GraphicsTyTank &&
		otherGraphics.GetId() == reflect.ValueOf(pAbsGraphics.sub).Interface().(*Bullet).TankId {
		return false
	}

	// 当前主动校验的是bullet时，不校验equip是否相交
	if pAbsGraphics.GraphicsTy == GraphicsTyBullet &&
		otherGraphics.GetGraphicsTy() == GraphicsTyEquip {
		return false
	}

	// 当前主动校验的是bullet时，不对tank发射的子弹集校验
	if pAbsGraphics.GraphicsTy == GraphicsTyBullet &&
		otherGraphics.GetGraphicsTy() == GraphicsTyBullet &&
		reflect.ValueOf(otherGraphics).Interface().(*Bullet).TankId == reflect.ValueOf(pAbsGraphics.sub).Interface().(*Bullet).TankId {
		return false
	}

	// 两个矩形相交机几种情况：images/rectangle_itersect.png
	// 重心距离在X轴和Y轴都小于两矩形的长或宽的一半之和

	width, height := pAbsGraphics.pImage.Size()
	//centerX := pAbsGraphics.Location.X + float64(width/2)
	//centerY := pAbsGraphics.Location.Y + float64(height/2)
	centerX := x + float64(width/2)
	if width%2 != 0 {
		centerX += 1
	}
	centerY := y + float64(height/2)
	if height%2 != 0 {
		centerY += 1
	}
	//log.Printf("center x: %v, y: %v\n", centerX, centerY)

	otherWidth, otherHeight := otherGraphics.GetImage().Size()
	otherCenterX := otherGraphics.GetLocation().X + float64(otherWidth/2)
	if otherWidth%2 != 0 {
		otherCenterX += 1
	}
	otherCenterY := otherGraphics.GetLocation().Y + float64(otherHeight/2)
	if otherHeight%2 != 0 {
		otherCenterY += 1
	}
	//log.Printf("otherCenter x: %v, y: %v\n", otherCenterX, otherCenterY)

	centerWidth := math.Abs(centerX - otherCenterX)
	centerHeight := math.Abs(centerY - otherCenterY)
	//log.Printf("center width: %v, height: %v\n", centerWidth, centerHeight)
	//log.Println()

	_width := (width + otherWidth) / 2
	if (width+otherWidth)%2 != 0 {
		_width += 1
	}
	_height := (height + otherHeight) / 2
	if (height+otherHeight)%2 != 0 {
		_height += 1
	}

	r := false
	if centerWidth <= float64(_width) && centerHeight <= float64(_height) {
		r = true
	}

	// 图形相交 & 当前主动校验的是tank，并且是当前端点的tank时，当前坦克捡起装备
	if r &&
		pAbsGraphics.GraphicsTy == GraphicsTyTank &&
		pAbsGraphics.Id == pApp.pGame.pTank.Id &&
		otherGraphics.GetGraphicsTy() == GraphicsTyEquip {
		pAbsGraphics.sub.(*Tank).WearEquip(otherGraphics.(*Equip))
		otherGraphics.(*Equip).Status = StatusTerm

		// 通知其它端点
		pApp.pEndpoint.SendGraphicsToAddrs(otherGraphics)

		return false
	}

	return r
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

	// 阻塞获取 chanel 中的 map
	graphicsMap := <-pApp.pGame.GraphicsMapChan

	// 再将 map 添加到 channel
	defer func() { pApp.pGame.GraphicsMapChan <- graphicsMap }()

	for _, value := range graphicsMap {
		if value.GetId() == pAbsGraphics.GetId() {
			continue
		}

		if pAbsGraphics.callIntersect(x, y, value) {
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
