// Package main contains HTTP server code.
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/vendelin8/card-fun/cmd/server/docs"
	"github.com/vendelin8/card-fun/internal/api"
	"github.com/vendelin8/card-fun/pkg/util"
)

func main() {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Put("/v1/create", api.CreateHandler)
	app.Get("/v1/open", api.OpenHandler)
	app.Post("/v1/draw", api.DrawHandler)

	backendAddr := util.GetEnv("BACKEND_URL", ":3000")
	if err := app.Listen(backendAddr); err != nil {
		log.Fatal(err)
	}
}
