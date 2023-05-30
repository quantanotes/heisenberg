package server

import (
	"encoding/json"
	"fmt"
	"github.com/quantanotes/heisenberg/core"
	"github.com/quantanotes/heisenberg/log"
	"github.com/quantanotes/heisenberg/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/keyauth/v2"
)

var apiKey = os.Getenv("HEISENBERG_API_KEY")
var masterKey = os.Getenv("HEISENBERG_MASTER_KEY")

type api struct {
	db *core.DB
}

func RunAPI(db *core.DB, host string) {
	api := api{db}
	app := fiber.New()
	app.Use(logger.New())
	app.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:X-API-Key",
		Validator: validateAPIKey,
	}))
	app.Post("/newcollection", api.handleNewCollection)
	app.Post("/deletecollection", api.handleDeleteCollection)
	app.Post("/put", api.handlePut)
	app.Post("/get", api.handleGet)
	app.Post("/delete", api.handleDelete)
	app.Post("/search", api.handleSearch)
	err := app.Listen(host)
	if err != nil {
		panic(err)
	}
}

func validateAPIKey(c *fiber.Ctx, key string) (bool, error) {
	if key == apiKey || key == masterKey {
		return true, nil
	}
	return false, fmt.Errorf("invalid api key")
}

func (a *api) handleNewCollection(c *fiber.Ctx) error {
	b := &struct {
		Name  string `json:"name"`
		Dim   uint   `json:"dim"`
		Space string `json:"space"`
	}{}
	if err := c.BodyParser(&b); err != nil {
		log.Error(err.Error(), nil)
		return c.Status(500).SendString(err.Error())
	}
	log.Trace(fmt.Sprintf("creating collection %s", b.Name), nil)
	space := utils.SpaceFromString(b.Space)
	err := a.db.NewCollection(b.Name, b.Dim, space)
	if err != nil {
		log.Error(err.Error(), nil)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(nil)
}

func (a *api) handleDeleteCollection(c *fiber.Ctx) error {
	b := &struct {
		Name string `json:"name"`
	}{}
	if err := c.BodyParser(&b); err != nil {
		log.Error(err.Error(), nil)
		return c.Status(500).SendString(err.Error())
	}
	err := a.db.DeleteCollection(b.Name)
	if err != nil {
		log.Error(err.Error(), nil)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(nil)
}

func (a *api) handlePut(c *fiber.Ctx) error {
	b := &struct {
		Collection string                 `json:"collection"`
		Key        string                 `json:"key"`
		Vector     []float32              `json:"vector"`
		Meta       map[string]interface{} `json:"meta"`
	}{}
	if err := c.BodyParser(&b); err != nil {
		log.Error(err.Error(), nil)
		return c.Status(500).SendString(err.Error())
	}
	err := a.db.Put(b.Collection, b.Key, b.Vector, b.Meta)
	if err != nil {
		log.Error(err.Error(), nil)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(nil)
}

func (a *api) handleGet(c *fiber.Ctx) error {
	b := &struct {
		Collection string `json:"collection"`
		Key        string `json:"key"`
	}{}
	if err := c.BodyParser(&b); err != nil {
		log.Error(err.Error(), nil)
		return c.Status(500).SendString(err.Error())
	}
	entry, err := a.db.Get(b.Collection, b.Key)
	if err != nil {
		log.Error(err.Error(), nil)
		return c.Status(500).SendString(err.Error())
	}
	data, err := json.Marshal(entry)
	if err != nil {
		log.Error(err.Error(), nil)
		return c.Status(500).SendString(err.Error())
	}
	return c.Send(data)
}

func (a *api) handleDelete(c *fiber.Ctx) error {
	b := &struct {
		Collection string `json:"collection"`
		Key        string `json:"key"`
	}{}
	if err := c.BodyParser(&b); err != nil {
		log.Error(err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err := a.db.Delete(b.Collection, b.Key)
	if err != nil {
		log.Error(err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(nil)
}

func (a *api) handleSearch(c *fiber.Ctx) error {
	b := &struct {
		Collection string    `json:"collection"`
		Query      []float32 `json:"query"`
		K          int       `json:"k"`
	}{}
	if err := c.BodyParser(&b); err != nil {
		log.Error(err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	res, err := a.db.Search(b.Collection, b.Query, b.K)
	if err != nil {
		log.Error(err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Error(err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Send(data)
}
