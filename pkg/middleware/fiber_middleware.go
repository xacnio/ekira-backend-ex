package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"io"
	"os"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(a *fiber.App) {
	logFile := "./logs/access.log"
	var logWriter io.Writer
	__logWriter, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("error opening access log file: %v", err)
		logWriter = os.Stdout
	} else {
		logWriter = __logWriter
	}

	a.Use(
		// Add CORS to each route.
		cors.New(),
		// Add simple logger.
		logger.New(
			logger.Config{
				Output:     logWriter,
				TimeFormat: "02-01-2006 15:04:05",
				TimeZone:   "Europe/Istanbul",
				Format:     "${time} ${status} - ${latency} ${method} ${path} ${ip} ${ua} ${error}\n",
			},
		),
	)
	a.Static("/", "./public")
}
