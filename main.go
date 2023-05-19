package main

import (
	"ChatDanBackend/bootstrap"
	_ "ChatDanBackend/docs"
	"ChatDanBackend/utils"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title          ChatDan Backend
// @version        0.0.1
// @description    ChatDan, a message box and 'biaobai' platform for Fudaners.
// @termsOfService https://swagger.io/terms/

// @tag.name User Module
// @tag.description 用户模块

// @tag.name MessageBox Module
// @tag.description 提问箱模块

// @tag.name Post Module
// @tag.description 帖子模块

// @tag.name Channel Module
// @tag.description 帖子回复 thread 模块

// @tag.name Wall Module
// @tag.description 表白墙模块

// @tag.name Division Module
// @tag.description 广场分区模块

// @tag.name Topic Module
// @tag.description 广场话题模块

// @tag.name Comment Module
// @tag.description 广场评论模块

// @tag.name Tag Module
// @tag.description 广场标签模块

// @tag.name Chat Module
// @tag.description 聊天模块

// @contact.name   JingYiJun
// @contact.url    https://www.jingyijun.xyz
// @contact.email  jingyijun3104@outlook.com

// @license.name  Apache 2.0
// @license.url   https://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

//go:generate go install github.com/swaggo/swag/cmd/swag@latest
//go:generate swag init
func main() {
	app := bootstrap.InitFiberApp()

	go func() {
		if innerErr := app.Listen("0.0.0.0:8000"); innerErr != nil {
			log.Println(innerErr)
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		utils.Logger.Error("app shutdown error", zap.Error(err))
	}

	_ = utils.Logger.Sync()
}
