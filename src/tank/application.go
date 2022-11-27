// run
//
// @author xiangqian
// @date 01:09 2022/11/27
package tank

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
)

// 图形Map
var graphicsMap map[string]Graphics

// 图形Map channel
var graphicsMapChan chan map[string]Graphics

func init() {
	graphicsMap = make(map[string]Graphics, 8)
	graphicsMapChan = make(chan map[string]Graphics, 1)
	graphicsMapChan <- graphicsMap
}

func clean0() {
	// 阻塞获取 chanel 中的 map
	graphicsMap := <-graphicsMapChan

	// 再将 map 添加到 channel
	defer func() { graphicsMapChan <- graphicsMap }()

	// 定义id切片（slice）
	var ids []string
	ids = nil
	index := 0
	for id, graphics := range graphicsMap {
		if !graphics.Alive() {
			if ids == nil {
				ids = make([]string, len(graphicsMap))
			}
			ids[index] = id
			index++
		}
	}

	if ids != nil {
		for i := 0; i < index; i++ {
			id := ids[i]
			delete(graphicsMap, id)
			fmt.Printf("delete %v\n", id)
		}
	}
}

func clean() {
	for {
		clean0()
	}
}

func draw0() {
	// 阻塞获取 chanel 中的 map
	graphicsMap := <-graphicsMapChan

	// 再将 map 添加到 channel
	defer func() { graphicsMapChan <- graphicsMap }()

	// 清除界面
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for _, graphics := range graphicsMap {
		if graphics.Alive() {
			graphics.Draw()
		}
	}

	// 刷新
	termbox.Flush()
}

func draw() {
	for {
		draw0()
		time.Sleep(30 * time.Millisecond)
	}
}

func addGraphics(graphics Graphics) {
	graphicsMap[graphics.Id()] = graphics
}

var tank *Tank

func graphics() {
	addGraphics(CreateBullet(&Location{10, 10}, DirectionRight, SpeedNormal))

	tank = CreateTank(Location{20, 20}, DirectionRight, SpeedNormal)
	addGraphics(tank)
}

func Run() {

	// 初始化 termbox
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	// 关闭 termbox
	defer termbox.Close()

	// 获取当前窗口宽度和高度
	// -→ x
	// ↓ y
	width, height := termbox.Size()
	fmt.Printf("width \t= %v\n", width)
	fmt.Printf("height \t= %v\n", height)

	//
	graphics()

	// 绘制
	go draw()

	// 清理
	go clean()

	// termbox事件(如,键盘按键) channel
	eventChan := make(chan termbox.Event)
	go func() {
		for {
			// 向 channel 添加轮询事件
			eventChan <- termbox.PollEvent()
		}
	}()

loop:
	for {
		select {
		case event := <-eventChan:
			// 如果是Key类型事件
			if event.Type == termbox.EventKey {
				switch event.Key {
				// Esc按键
				case termbox.KeyEsc:
					break loop

				// 向上键箭头按键
				case termbox.KeyArrowUp:
					tank.Move(DirectionUp)

				// 向下键箭头按键
				case termbox.KeyArrowDown:
					tank.Move(DirectionDown)

				// 向左键箭头按键
				case termbox.KeyArrowLeft:
					tank.Move(DirectionLeft)

				// 向右键箭头按键
				case termbox.KeyArrowRight:
					tank.Move(DirectionRight)
				}
			}
		default:
			time.Sleep(2 * time.Millisecond)
		}
	}

}
