package mockservice

import (
	"github.com/gofiber/fiber/v2"
)

func (a *API) registerRestRoutes(router fiber.Router) {
	r := fiber.New()
	r.Post("/", a.createRestService)
	r.Put("/:id", a.updateRestService)
	r.Delete("/:id", a.deleteRestService)
	router.Mount("/rest", r)
}

// createRestService godoc
// @ID createRestService
// @Summary "create rest service"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} Create
func (a *API) createRestService(ctx *fiber.Ctx) error {

}

func (a *API) updateRestService(ctx *fiber.Ctx) error {

}

func (a *API) deleteRestService(ctx *fiber.Ctx) error {

}
