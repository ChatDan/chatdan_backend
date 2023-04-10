package fiberx

import (
	"ChatDanBackend/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"time"
)

func registerMiddlewares(app *fiber.App) *fiber.App {
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(MyLogger)
	app.Use(pprof.New())
	return app
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
	common.Logger.Info("http log", output...)
	return nil
}
