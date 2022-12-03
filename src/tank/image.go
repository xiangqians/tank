// image
// @author xiangqian
// @date 22:31 2022/11/29
package tank

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

type Image struct {
	pDefaultImage     *ebiten.Image // default
	pTankUpImage      *ebiten.Image // tank
	pTankDownImage    *ebiten.Image
	pTankLeftImage    *ebiten.Image
	pTankRightImage   *ebiten.Image
	pBulletUpImage    *ebiten.Image // bullet
	pBulletDownImage  *ebiten.Image
	pBulletLeftImage  *ebiten.Image
	pBulletRightImage *ebiten.Image
}

func (pImage *Image) Init() {
	// default
	pImage.pDefaultImage = NewEbitenImage("images/default.png")

	// tank
	pImage.pTankUpImage = NewEbitenImage("images/tankU.gif")
	pImage.pTankDownImage = NewEbitenImage("images/tankD.gif")
	pImage.pTankLeftImage = NewEbitenImage("images/tankL.gif")
	pImage.pTankRightImage = NewEbitenImage("images/tankR.gif")

	// bullet
	pImage.pBulletUpImage = NewEbitenImage("images/bulletU.gif")
	pImage.pBulletDownImage = NewEbitenImage("images/bulletD.gif")
	pImage.pBulletLeftImage = NewEbitenImage("images/bulletL.gif")
	pImage.pBulletRightImage = NewEbitenImage("images/bulletR.gif")
}

func NewEbitenImage(name string) *ebiten.Image {
	pFile, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	// image: unknown format
	// 需要导入 import _ "image/gif"，image包不知道怎么Decode图片，需要导入 "image/gif" 去解码 gif 图片。
	img, _, err := image.Decode(pFile)
	if err != nil {
		log.Fatal(err)
	}

	pEbitenImage, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	return pEbitenImage
}
