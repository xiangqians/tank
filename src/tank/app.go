// Application
// @author xiangqian
// @date 22:22 2022/11/29
package tank

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"os"
	"time"
)

// App步骤
type AppStep int8

const (
	AppStepReg  AppStep = iota // 注册界面
	AppStepGame                // 游戏界面
)

const (
	// 1280 * 720
	screenWidth  = 1280
	screenHeight = 720
)

var pApp *App

func init() {
	initLog()
}

func initLog() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 创建日志文件夹，如果不存在的话
	logDir := "./logs"
	fileInfo, err := os.Stat(logDir)
	if err != nil || !fileInfo.IsDir() {
		err = os.Mkdir(logDir, 0666)
		if err != nil {
			panic(err)
		}
	}

	// 创建日志文件
	pLogFile, err := os.OpenFile(logDir+"/tank.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	// 设置日志输出
	log.SetOutput(pLogFile)
}

type App struct {
	pEndpoint              *Endpoint // 端点
	pReg                   *Reg      // 注册界面
	pGame                  *Game     // 游戏界面
	appStep                AppStep   // App步骤
	prevKeyPressedUnixNano int64     // 上一次按键 unix nano
	curKeyPressedUnixNano  int64     // 当前按键 unix nano
}

func (pApp *App) init() {
	pApp.pEndpoint = &Endpoint{}
	pApp.pReg = &Reg{}
	pApp.pGame = &Game{}

	// 异步监听端点
	go pApp.pEndpoint.Listen()

	// init
	pApp.pReg.init()
	pApp.pGame.init()
}

func (pApp *App) Update(screen *ebiten.Image) error {
	switch pApp.appStep {
	case AppStepReg:
		return pApp.pReg.Update(screen)
	case AppStepGame:
		return pApp.pGame.Update(screen)
	}
	return nil
}

func (pApp *App) IsKeyPressed(key ebiten.Key) bool {
	if ebiten.IsKeyPressed(key) {
		pApp.curKeyPressedUnixNano = time.Now().UnixNano()
		result := pApp.curKeyPressedUnixNano-pApp.prevKeyPressedUnixNano >= 10
		pApp.prevKeyPressedUnixNano = pApp.curKeyPressedUnixNano
		return result
	}
	return false
}

func (pApp *App) Draw(screen *ebiten.Image) {
	switch pApp.appStep {
	case AppStepReg:
		pApp.pReg.Draw(screen)
	case AppStepGame:
		pApp.pGame.Draw(screen)
	}
}

func (pApp *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Run() {
	// app
	pApp = &App{}
	pApp.init()

	// RUN
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tank")
	if err := ebiten.RunGame(pApp); err != nil {
		log.Fatal(err)
	}
}
