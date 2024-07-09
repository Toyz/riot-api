package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"riot-api/internal/henrik"
	"time"
)

const cacheTime = 4 * time.Minute

type ValorantController struct {
	session *session.Store
	redis   fiber.Storage
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

	key := c.createCacheKey(region, name, tag, "rank")
	if c.redis != nil {
		if data, err := c.redis.Get(key); err == nil && len(data) > 0 {
			ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			ctx.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=%v", cacheTime.Seconds()))
			return ctx.Send(data)
		}
	}

	rank, err := henrik.GetRank(region, name, tag)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if c.redis != nil {
		data, err := json.Marshal(rank)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := c.redis.Set(key, data, cacheTime); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	ctx.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=%v", cacheTime.Seconds()))
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

	key := c.createCacheKey(region, name, tag, "matches")
	if c.redis != nil {
		if data, err := c.redis.Get(key); err == nil && len(data) > 0 {
			ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			ctx.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=%v", cacheTime.Seconds()))
			return ctx.Send(data)
		}
	}

	matches, err := henrik.GetMatches(region, name, tag)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if c.redis != nil {
		data, err := json.Marshal(matches)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := c.redis.Set(key, data, cacheTime); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	ctx.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=%v", cacheTime.Seconds()))
	return ctx.JSON(matches)
}

func (*ValorantController) createCacheKey(region, name, tag, lookUpTType string) string {
	key := region + "-" + name + "-" + tag + "-" + lookUpTType
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%x_valorant", hash)
}

func NewValorantController(app fiber.Router, session *session.Store, redis fiber.Storage) {
	api := &ValorantController{session: session, redis: redis}

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
