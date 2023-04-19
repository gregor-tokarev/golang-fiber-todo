package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"goapi/models"
	"goapi/router"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	models.InitDB()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowMethods:     "*",
	}))

	api := app.Group("/api/v1")

	router.InitUserRoutes(api)
	router.InitTaskRouter(api)

	log.Fatal(app.Listen(":3000"))
}
