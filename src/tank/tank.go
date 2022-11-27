// tank
// 坦克
//
// @author xiangqian
// @date 10:58 2022/11/27
package tank

import "fmt"

type Tank struct {
	body []Location // 坦克体
	*AbsGraphics
}

func (tank *Tank) Draw() error {
	fmt.Printf("draw tank\n")
	tank.alive = false

	// 坦克
	//termbox.SetBg(20, 10, termbox.ColorRed)
	//termbox.SetBg(19, 11, termbox.ColorRed)
	//termbox.SetBg(21, 11, termbox.ColorRed)

	return nil
}
