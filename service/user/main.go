package main

import (
	"ChatDanBackend/common/fiberx"
	"ChatDanBackend/service/user/api"
	_ "ChatDanBackend/service/user/config"
)

func main() {
	app := fiberx.NewFiberApp("Message Box")
	api.RegisterRoutes(app)
	fiberx.AppListen(app)
}
