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
	pTankUpImg    *ebiten.Image
	pTankDownImg  *ebiten.Image
	pTankLeftImg  *ebiten.Image
	pTankRightImg *ebiten.Image
)

func init() {
	pTankUpImg = newImage("images/tankU.gif")
	pTankDownImg = newImage("images/tankD.gif")
	pTankLeftImg = newImage("images/tankL.gif")
	pTankRightImg = newImage("images/tankR.gif")
}

func newImage(name string) *ebiten.Image {
	pFile, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	// image: unknown format
	// 需要导入 import _ "image/gif"，image包不知道怎么Decode图片，需要导入 "image/gif" 去解码 gif 图片。
	pImg, _, err := image.Decode(pFile)
	if err != nil {
		log.Fatal(err)
	}

	pEbitenImg, err := ebiten.NewImageFromImage(pImg, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	return pEbitenImg
}
