package main

import (
	"ChatDanBackend/common/fiberx"
	"ChatDanBackend/service/user/api"
	_ "ChatDanBackend/service/user/config"
	_ "ChatDanBackend/service/user/docs"
)

// @title           User Microservice
// @version         0.0.1
// @description     This is a user service for ChatDan.
// @termsOfService  https://swagger.io/terms/

// @contact.name   JingYiJun
// @contact.url    https://danxi.fduhole.com
// @contact.email  jingyijun@fduhole.com

// @license.name  Apache 2.0
// @license.url   https://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	app := fiberx.NewFiberApp("Message Box")
	api.RegisterRoutes(app)
	fiberx.AppListen(app)
}
