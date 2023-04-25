package fiberx

import (
	"ChatDanBackend/common/cachex"
	"ChatDanBackend/common/configx"
	"ChatDanBackend/common/gormx"
	"ChatDanBackend/common/loggerx"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type AppOptions struct {
	AppName        string
	CustomConfig   any
	Models         []any
	RegisterRoutes func(app *fiber.App)
}

// NewFiberApp creates a new fiber app and listen on 0.0.0.0:8000
func NewFiberApp(options AppOptions) {
	// bootstrap
	configx.InitConfig(options.CustomConfig)
	gormx.InitDB(options.Models...)
	cachex.InitCache()

	// new fiber app
	app := fiber.New(fiber.Config{
		AppName:      options.AppName,
		ErrorHandler: MyErrorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
	registerMiddlewares(app)
	registerRoutes(app)
	if options.RegisterRoutes != nil {
		options.RegisterRoutes(app)
	}

	// listen
	AppListen(app)
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
		loggerx.Logger.Error("app shutdown error", zap.Error(err))
	}

	// sync logger
	err = loggerx.Logger.Sync()
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
