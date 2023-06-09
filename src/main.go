package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"goapi/src/config"
	"goapi/src/models"
	router2 "goapi/src/router"
	"log"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	models.InitDB()
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(recover2.New())
	app.Use(logger.New())

	app.Get("/monitor", monitor.New())

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowMethods:     "*",
	}))
	app.Use(limiter.New(limiter.Config{
		Next: func(ctx *fiber.Ctx) bool {
			return ctx.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 30 * time.Second,
		LimitReached: func(ctx *fiber.Ctx) error {
			return ctx.Status(429).JSON(fiber.Map{
				"message": "Too many requests",
			})
		},
	}))
	app.Use(cache.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	api := app.Group(config.Cfg.ServerPrefix)

	router2.InitUserRoutes(api)
	router2.InitTaskRouter(api)
	router2.InitTagRouter(api)

	log.Fatal(app.Listen(":3000"))
}
