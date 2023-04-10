package main

import (
	"ChatDanBackend/common/fiberx"
	_ "ChatDanBackend/service/message_box/config"
)

func main() {
	app := fiberx.NewFiberApp("Message Box")
	// todo: register routes
	fiberx.AppListen(app)
}
