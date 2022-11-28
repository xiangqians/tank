// information bar
// 信息栏
// @author xiangqian
// @date 22:15 2022/11/28
package tank

import "github.com/nsf/termbox-go"

type InfoBar struct {
	*AbsGraphics
	width int
}

func CreateInfoBar(width int) *InfoBar {
	return &InfoBar{
		AbsGraphics: &AbsGraphics{
			id: "infoBar_" + Uuid(),
		},
		width: width,
	}
}

func (infoBar *InfoBar) Draw() error {
	if err := infoBar.AbsGraphics.Draw(); err != nil {
		return err
	}

	_, height := termbox.Size()

	for y := 0; y <= height; y++ {
		termbox.SetBg(infoBar.width, y, termbox.ColorWhite)
	}

	return nil
}
