// bullet
// 子弹
//
// @author xiangqian
// @date 11:01 2022/11/27
package tank

import (
	"github.com/hajimehoshi/ebiten"
	"sync"
	"time"
)

// 互斥锁
var bulletLock sync.Mutex

type Bullet struct {
	*AbsGraphics
	tankId string // 子弹所属坦克id
}

func CreateBullet(pTank *Tank, speed Speed) *Bullet {
	// 让子弹坐标从坦克中心发出
	pLocation := pTank.pLocation
	width, height := pTank.pImg.Size()
	//log.Printf("width = %v, height = %v\n", width, height)
	var location Location
	switch pTank.direction {
	case DirectionUp:
		location = Location{pLocation.x + float64(width/2) - 4.51392, pLocation.y - 8}
	case DirectionDown:
		location = Location{pLocation.x + float64(width/2) - 5.8, pLocation.y + float64(height)}
	case DirectionLeft:
		location = Location{pLocation.x - 8, pLocation.y + float64(height/2) - 2.91392}
	case DirectionRight:
		location = Location{pLocation.x + float64(width), pLocation.y + float64(height/2) - 2.91392}
	}
	//log.Printf("location = %v\n", location)

	pBullet := &Bullet{
		AbsGraphics: CreateAbsGraphics(location, pTank.direction, speed),
		tankId:      pTank.id,
	}
	pBullet.Init(pBullet)
	return pBullet
}

func (bullet *Bullet) IsOutOfBounds(x, y float64) bool {
	r := bullet.AbsGraphics.IsOutOfBounds(x, y)
	if r {
		bullet.status = StatusTerm
	}
	return r
}

func (bullet *Bullet) Draw(screen *ebiten.Image) error {
	if err := bullet.AbsGraphics.Draw(screen); err != nil {
		return err
	}

	// 如果是当前用户所发射的子弹，那么由当前用户轮询设置子弹位置
	if bullet.tankId == pTank.id && bullet.status == StatusNew {
		func() {
			// 加锁
			bulletLock.Lock()

			if bullet.status == StatusNew {
				go bullet.Run()
				bullet.status = StatusRun
			}

			// 释放锁
			defer bulletLock.Unlock()
		}()
	}

	return nil
}

func (bullet *Bullet) Run() {
	for {
		if bullet.status != StatusRun {
			break
		}

		bullet.Move(bullet.direction)
		switch bullet.speed {
		case SpeedSlow:
			time.Sleep(100 * time.Millisecond)

		case SpeedNormal:
			time.Sleep(50 * time.Millisecond)

		case SpeedFast:
			time.Sleep(10 * time.Millisecond)
		}
	}
	//log.Printf("%v Term\n", bullet.id)
}

func (bullet *Bullet) UpImg() *ebiten.Image {
	return pBulletUpImg
}

func (bullet *Bullet) DownImg() *ebiten.Image {
	return pBulletDownImg
}

func (bullet *Bullet) LeftImg() *ebiten.Image {
	return pBulletLeftImg
}

func (bullet *Bullet) RightImg() *ebiten.Image {
	return pBulletRightImg
}
