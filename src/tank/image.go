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

var (
	pDefaultImg     *ebiten.Image // default
	pTankUpImg      *ebiten.Image // tank
	pTankDownImg    *ebiten.Image
	pTankLeftImg    *ebiten.Image
	pTankRightImg   *ebiten.Image
	pBulletUpImg    *ebiten.Image // bullet
	pBulletDownImg  *ebiten.Image
	pBulletLeftImg  *ebiten.Image
	pBulletRightImg *ebiten.Image
)

func init() {
	// default
	pDefaultImg = NewEbitenImage("images/default.png")

	// tank
	pTankUpImg = NewEbitenImage("images/tankU.gif")
	pTankDownImg = NewEbitenImage("images/tankD.gif")
	pTankLeftImg = NewEbitenImage("images/tankL.gif")
	pTankRightImg = NewEbitenImage("images/tankR.gif")

	// bullet
	pBulletUpImg = NewEbitenImage("images/bulletU.gif")
	pBulletDownImg = NewEbitenImage("images/bulletD.gif")
	pBulletLeftImg = NewEbitenImage("images/bulletL.gif")
	pBulletRightImg = NewEbitenImage("images/bulletR.gif")
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

	pEbitenImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	return pEbitenImg
}
