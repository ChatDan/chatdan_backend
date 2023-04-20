package main

import (
	"ChatDanBackend/common/fiberx"
	"ChatDanBackend/service/wall/model"
)

func main() {
	fiberx.NewFiberApp(fiberx.AppOptions{
		AppName:        "Wall Microservice",
		CustomConfig:   nil,
		Models:         []any{model.Wall{}},
		RegisterRoutes: nil,
	})
}
