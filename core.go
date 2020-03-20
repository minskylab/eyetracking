package main

import (
	"time"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var storeInfo *StorageInfo
var lastFetch time.Time

func fetchStoreState() {
	log.Info("fetching current images on store")
	var err error
	storeInfo, err = storageInfo()
	if err != nil {
		log.Error(errors.Wrap(err, "error at try to fetch store state"))
	}
	lastFetch = time.Now()
}

func init() {
	minioClientInit()
	fetchStoreState()
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://minsky.cc", "projects.minsky.cc"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"*"},
	}))

	apiv1 := app.Group("/api/v1")

	apiv1.Get("/info", func(c *fiber.Ctx) {
		if time.Now().Sub(lastFetch) > 5*time.Second {
			fetchStoreState()
		}

		if err := c.JSON(storeInfo); err != nil {
			log.Error(errors.Wrap(err, "error at json serializing reponse"))
		}
	})

	if err := app.Listen(8080); err != nil {
		err = errors.Wrap(err, "error at listen service")
		log.Error(err)
		panic(err)
	}
}
