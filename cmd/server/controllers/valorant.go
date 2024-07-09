package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"riot-api/internal/henrik"
)

type ValorantController struct {
	session *session.Store
}

func (c *ValorantController) getCurrentRank(ctx *fiber.Ctx) error {
	region := ctx.Params("region", "na")
	name := ctx.Params("name", "")
	tag := ctx.Params("tag", "")

	if name == "" || tag == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name and tag are required",
		})
	}

	rank, err := henrik.GetRank(region, name, tag)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(rank)
}

func (c *ValorantController) getMatches(ctx *fiber.Ctx) error {
	region := ctx.Params("region", "na")
	name := ctx.Params("name", "")
	tag := ctx.Params("tag", "")

	if name == "" || tag == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name and tag are required",
		})
	}

	matches, err := henrik.GetMatches(region, name, tag)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(matches)
}

func NewValorantController(app fiber.Router, session *session.Store) {
	api := &ValorantController{session: session}

	cors := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,HEAD,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	})
	app.Use(cors)

	party := app.Group("/valorant/:region")
	party.Get("/rank/:name/:tag", api.getCurrentRank)
	party.Get("/matches/:name/:tag", api.getMatches)
}
