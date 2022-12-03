// bullet
// 子弹
//
// @author xiangqian
// @date 11:01 2022/11/27
package tank

import (
	"github.com/hajimehoshi/ebiten"
	"log"
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
	pLocation := pTank.Location
	width, height := pTank.pImg.Size()
	//log.Printf("width = %v, height = %v\n", width, height)
	var location Location
	switch pTank.Direction {
	case DirectionUp:
		location = Location{pLocation.X + float64(width/2) - 4.51392, pLocation.Y - 8}
	case DirectionDown:
		location = Location{pLocation.X + float64(width/2) - 5.8, pLocation.Y + float64(height)}
	case DirectionLeft:
		location = Location{pLocation.X - 8, pLocation.Y + float64(height/2) - 2.91392}
	case DirectionRight:
		location = Location{pLocation.X + float64(width), pLocation.Y + float64(height/2) - 2.91392}
	}
	//log.Printf("location = %v\n", location)

	pBullet := &Bullet{
		AbsGraphics: CreateAbsGraphics(GraphicsTyBullet, location, pTank.Direction, speed),
		tankId:      pTank.Id,
	}
	pBullet.Init(pBullet)
	return pBullet
}

func (pBullet *Bullet) IsOutOfBounds(x, y float64) bool {
	r := pBullet.AbsGraphics.IsOutOfBounds(x, y)
	if r {
		pBullet.Status = StatusTerm
	}
	return r
}

func (pBullet *Bullet) Draw(screen *ebiten.Image) error {
	if err := pBullet.AbsGraphics.Draw(screen); err != nil {
		return err
	}

	// 如果是当前用户所发射的子弹，那么由当前用户轮询设置子弹位置
	if pBullet.tankId == pApp.pGame.pTank.Id && pBullet.Status == StatusNew {
		func() {
			// 加锁
			bulletLock.Lock()

			if pBullet.Status == StatusNew {
				go pBullet.Run()
				pBullet.Status = StatusRun
			}

			// 释放锁
			defer bulletLock.Unlock()
		}()
	}

	return nil
}

func (pBullet *Bullet) Run() {
	for {
		if pBullet.Status != StatusRun {
			break
		}

		pBullet.Move(pBullet.Direction)

		buf, _ := Serialize(pBullet)
		log.Printf("%v\n", string(buf))

		switch pBullet.Speed {
		case SpeedSlow:
			time.Sleep(100 * time.Millisecond)

		case SpeedNormal:
			time.Sleep(50 * time.Millisecond)

		case SpeedFast:
			time.Sleep(10 * time.Millisecond)
		}
	}
	//log.Printf("%v Term\n", pBullet.id)
}

func (pBullet *Bullet) UpImg() *ebiten.Image {
	return pBulletUpImg
}

func (pBullet *Bullet) DownImg() *ebiten.Image {
	return pBulletDownImg
}

func (pBullet *Bullet) LeftImg() *ebiten.Image {
	return pBulletLeftImg
}

func (pBullet *Bullet) RightImg() *ebiten.Image {
	return pBulletRightImg
}
