package mockservice

import (
	"configurator/internal/api"
	"configurator/internal/entity"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type API struct {
	log     *zap.SugaredLogger
	useCase IUseCase
}

func New(log *zap.SugaredLogger, useCase IUseCase) *API {
	return &API{
		log:     log,
		useCase: useCase,
	}
}

func (a *API) Register(router fiber.Router, authenticator api.Authenticator, middlewares ...fiber.Handler) {
	r := fiber.New()
	r.Use("/", authenticator.AuthMiddleware())
	a.registerRestRoutes(r)
	r.Get("/", a.listServices)
	router.Mount("/mockservices", r)
}

// listServices godoc
// @ID listServices
// @Summary "list all the services in project"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} listServiceResponse
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to list all services"
// @Router /mockservices [get]
func (a *API) listServices(ctx *fiber.Ctx) error {
	projectId, err := api.ParseQueryParam[entity.BigIntPK](ctx, "projectId")
	if err != nil {
		return err
	}
	result, err := a.useCase.ListServices(ctx.Context(), projectId)
	if err != nil {
		return err
	}

	response := listServiceRespByUseCaseResult(result)
	return ctx.JSON(response)
}
