package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/redis/v2"
	"github.com/mikhail-bigun/fiberlogrus"
	log "github.com/sirupsen/logrus"
	"riot-api/cmd/server/controllers"
	"time"

	"github.com/spf13/viper"
	_ "riot-api/cmd"
)

func main() {
	storage := redis.New(redis.Config{
		URL:   viper.GetString("connection_strings.redis"),
		Reset: false,
	})

	app := fiber.New(fiber.Config{
		//Views:             engine,
		//ViewsLayout:       "layouts/layout",
		EnablePrintRoutes: true,
	})

	app.Use(recover.New(
		recover.Config{
			EnableStackTrace: true,
		}))
	app.Use(fiberlogrus.New(fiberlogrus.Config{
		Logger: log.StandardLogger(),
		Tags: []string{
			fiberlogrus.TagMethod,
			fiberlogrus.TagStatus,
			fiberlogrus.TagPath,
			fiberlogrus.TagLatency,
			fiberlogrus.TagIP,
			fiberlogrus.TagUA,
			fiberlogrus.AttachKeyTag(fiberlogrus.TagLocals, "requestid"),
		},
	}))
	app.Use(requestid.New())

	sess := session.New(session.Config{
		Expiration:     60 * time.Minute,
		Storage:        storage,
		CookieSameSite: "Lax",
		KeyGenerator:   utils.UUIDv4,
	})

	apiRoute := app.Group("/api")
	controllers.NewValorantController(apiRoute, sess)

	if err := app.Listen(viper.GetString("webserver.address")); err != nil {
		log.Fatal(err)
	}
}
