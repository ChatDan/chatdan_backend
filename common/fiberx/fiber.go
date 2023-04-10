package fiberx

import (
	"ChatDanBackend/common"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func NewFiberApp(appName string) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      appName,
		ErrorHandler: MyErrorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
	registerRoutes(app)
	return registerMiddlewares(app)
}

func AppListen(app *fiber.App) {
	go func() {
		err := app.Listen("0.0.0.0:8000")
		if err != nil {
			log.Fatal(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		common.Logger.Error("app shutdown error", zap.Error(err))
	}

	// sync logger
	err = common.Logger.Sync()
	if err != nil {
		log.Println(err)
	}
}

func registerRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api")
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", swagger.HandlerDefault)
}
