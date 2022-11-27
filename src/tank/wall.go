// wall
// 墙
//
// @author xiangqian
// @date 11:12 2022/11/27
package tank

import "fmt"

type Wall struct {
	body Location // 墙体
	*AbsGraphics
}

func (wall *Wall) Draw() error {
	fmt.Printf("draw wall\n")
	return nil
}
