package bootstrap

import (
	"ChatDanBackend/apis"
	"ChatDanBackend/config"
	"ChatDanBackend/models"
	"ChatDanBackend/utils"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"time"
)

func InitFiberApp() *fiber.App {
	config.InitConfig()
	models.InitDB()

	app := fiber.New(fiber.Config{
		AppName:               config.Config.AppName,
		ErrorHandler:          utils.MyErrorHandler,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	registerMiddlewares(app)
	apis.RegisterRoutes(app)

	return app
}

func registerMiddlewares(app *fiber.App) {
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(compress.New())
	app.Use(MyLogger)
	app.Use(pprof.New())
	if config.Config.Debug {
		app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	}
}

func MyLogger(c *fiber.Ctx) error {
	startTime := time.Now()
	chainErr := c.Next()

	if chainErr != nil {
		if err := c.App().ErrorHandler(c, chainErr); err != nil {
			_ = c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	latency := time.Since(startTime).Milliseconds()
	userID, ok := c.Locals("user_id").(int)
	output := []zap.Field{
		zap.Int("status_code", c.Response().StatusCode()),
		zap.String("method", c.Method()),
		zap.String("origin_url", c.OriginalURL()),
		zap.String("remote_ip", c.Get("X-Real-IP")),
		zap.Int64("latency", latency),
	}
	if ok {
		output = append(output, zap.Int("user_id", userID))
	}
	if chainErr != nil {
		output = append(output, zap.Error(chainErr))
	}
	utils.Logger.Info("http log", output...)
	return nil
}
