package main

import (
	"ChatDanBackend/common/fiberx"
	"ChatDanBackend/service/message_box/api"
	_ "ChatDanBackend/service/message_box/config"
	_ "ChatDanBackend/service/message_box/docs"
)

// @title           MessageBox Microservice
// @version         0.0.1
// @description     This is a message box service for ChatDan.
// @termsOfService  https://swagger.io/terms/

// @contact.name   PinappleUnderTheSea
// @contact.url    https://github.com/PinappleUnderTheSea
// @contact.email  hastaluego@fduhole.com

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
