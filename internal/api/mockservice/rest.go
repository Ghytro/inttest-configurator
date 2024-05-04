package mockservice

import (
	"configurator/internal/api"
	"configurator/internal/entity"
	"configurator/internal/usecase/mockservice"

	"github.com/gofiber/fiber/v2"
)

func (a *API) registerRestRoutes(router fiber.Router) {
	r := fiber.New()

	r.Post("/", a.createRestService).
		Put("/:id", a.updateRestService).
		Delete("/:id", a.deleteRestService)

	handGroup := r.Group("/:serviceId/handlers").
		Get("/", a.listRestHandlers).
		Post("/", a.addRestHandler).
		Put("/:id", a.updateRestHandler).
		Delete("/:id", a.deleteRestHandler)

	handGroup.Group("/:handlerId").
		Get("/behaviors", a.listRestBehaviors).
		Post("/addRestStubBehavior", a.addRestStubBehavior).
		Post("/addRestMockBehavior", a.addRestMockBehavior).
		Put("/updateRestStubBehavior/:id", a.updateRestStubBehavior).
		Put("/updateRestMockBehavior/:id", a.updateRestMockBehavior).
		Patch("/behaviors/:id/move", a.moveRestBehaviorPriority).
		Delete("/behaviors/:id", a.deleteRestBehavior)

	router.Mount("/rest", r)
}

// createRestService godoc
// @ID createRestService
// @Summary "create rest service"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} mockservice.CreateRestServiceResult
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to list all services"
// @Param form body mockservice.CreateRestServiceReq true "body form"
// @Router /mockservices/rest [post]
func (a *API) createRestService(ctx *fiber.Ctx) error {
	projectId, err := api.ParseQueryParam[entity.BigIntPK](ctx, "projectId")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.CreateRestServiceReq](ctx)
	if err != nil {
		return err
	}

	result, err := a.useCase.CreateRestService(ctx.Context(), projectId, form)
	if err != nil {
		return err
	}
	return ctx.JSON(result)
}

// updateRestService godoc
// @ID updateRestService
// @Summary "update rest service"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} mockservice.RestServiceUpdateResult
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to list all services"
// @Param id path string true "id of service to update"
// @Param form body mockservice.RestServiceUpdateForm true "body form"
// @Router /mockservices/rest/{id} [put]
func (a *API) updateRestService(ctx *fiber.Ctx) error {
	serviceId, err := api.ParseUrlParamsId[string](ctx, "id")
	if err != nil {
		return err
	}
	projectId, err := api.ParseQueryParam[entity.BigIntPK](ctx, "projectId")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.RestServiceUpdateForm](ctx)
	if err != nil {
		return err
	}
	result, err := a.useCase.UpdateRestService(
		ctx.Context(),
		mockservice.RestServiceIdentifier{
			ProjectId: projectId,
			ServiceId: serviceId,
		},
		form,
	)
	if err != nil {
		return err
	}
	return ctx.JSON(result)
}

// deleteRestService godoc
// @ID deleteRestService
// @Summary delete rest service
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param id path string true "id of a service to delete"
// @Param projectId query number true "id of a project where service is located"
// @Router /mockservices/rest/{id} [delete]
func (a *API) deleteRestService(ctx *fiber.Ctx) error {
	serviceId, err := api.ParseUrlParamsId[string](ctx, "id")
	if err != nil {
		return err
	}
	projectId, err := api.ParseQueryParam[entity.BigIntPK](ctx, "projectId")
	if err != nil {
		return err
	}
	err = a.useCase.DeleteRestService(ctx.Context(), mockservice.RestServiceIdentifier{
		ProjectId: projectId,
		ServiceId: serviceId,
	})
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// listRestHandlers godoc
// @ID listRestHandlers
// @Summary list all the rest handlers available
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {array} mockservice.ListRestRoutesResult
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Route /mockservices/rest/{serviceId}/handlers [get]
func (a *API) listRestHandlers(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}

	result, err := a.useCase.ListRestHandlers(ctx.Context(), serviceId)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

// addRestHandler godoc
// @ID addRestHandler
// @Summary add new rest handler
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Param form body mockservice.CreateRestHandlerForm true "body model"
// @Route /mockservices/rest/{serviceId}/handlers [post]
func (a *API) addRestHandler(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.CreateRestHandlerForm](ctx)
	if err != nil {
		return err
	}
	err = a.useCase.AddRestHandler(ctx.Context(), serviceId, form)
	if err != nil {
		return err
	}

	return ctx.JSON(api.OK{})
}

// updateRestHandler godoc
// @ID updateRestHandler
// @Summary update existing rest handler
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Param id path number true "id of a rest handler"
// @Param form body mockservice.UpdateRestHandlerForm true "body form"
// @Route /mockservices/rest/{serviceId}/handlers/{id} [put]
func (a *API) updateRestHandler(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	id, err := api.ParseUrlParamsId[int](ctx, "id")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.UpdateRestHandlerForm](ctx)
	if err != nil {
		return err
	}
	err = a.useCase.UpdateRestHandler(ctx.Context(), serviceId, id, form)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// deleteRestHandler godoc
// @ID deleteRestHandler
// @Summary delete existing rest handler
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Route /mockservices/rest/{serviceId}/handlers/{id} [delete]
func (a *API) deleteRestHandler(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	id, err := api.ParseUrlParamsId[int](ctx, "id")
	if err != nil {
		return err
	}
	err = a.useCase.DeleteRestHandler(ctx.Context(), serviceId, id)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// addRestStubBehavior godoc
// @ID addRestStubBehavior
// @Summary add new stub behavior to rest route
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Param handlerId path number true "id of rest handler"
// @Param form body mockservice.CreateRestStubBehaviorForm true "body model"
// @Route /mockservices/rest/{serviceId}/handlers/{handlerId}/addRestStubBehavior [post]
func (a *API) addRestStubBehavior(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	handlerId, err := api.ParseUrlParamsId[int](ctx, "handlerId")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.CreateRestStubBehaviorForm](ctx)
	if err != nil {
		return err
	}
	err = a.useCase.AddRestStubBehavior(ctx.Context(), serviceId, handlerId, form)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// addRestMockBehavior godoc
// @ID addRestMockBehavior
// @Summary add new mock behavior to rest route
// @Tags mockservices
// @Accept json
// @Produce json
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Param handlerId path number true "id of rest handler"
// @Param form body mockservice.CreateRestMockBehaviorForm true "body model"
// @Route /mockservices/rest/{serviceId}/handlers/{handlerId}/addRestMockBehavior [post]
func (a *API) addRestMockBehavior(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	handlerId, err := api.ParseUrlParamsId[int](ctx, "handlerId")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.CreateRestMockBehaviorForm](ctx)
	if err != nil {
		return err
	}
	err = a.useCase.AddRestMockBehavior(ctx.Context(), serviceId, handlerId, form)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// listRestBehaviors godoc
// @ID listRestBehaviors
// @Summary list all available behaviors of rest handler
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} mockservice.ListRestBehaviorResult
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Param handlerId path number true "id of rest handler"
// @Route /mockservices/rest/{serviceId}/handlers/{handlerId}/behaviors [get]
func (a *API) listRestBehaviors(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	handlerId, err := api.ParseUrlParamsId[int](ctx, "handlerId")
	if err != nil {
		return err
	}
	result, err := a.useCase.ListRestBehaviors(ctx.Context(), serviceId, handlerId)
	if err != nil {
		return err
	}
	return ctx.JSON(result)
}

// updateRestStubBehavior godoc
// @ID updateRestStubBehavior
// @Summary update existing stub behavior
// @Tags mockservices
// @Accept json
// @Produce json
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Param handlerId path number true "id of rest handler"
// @Param id path number true "id of behavior to update"
// @Route /mockservices/rest/{serviceId}/handlers/{handlerId}/updateRestStubBehavior/{id} [put]
func (a *API) updateRestStubBehavior(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	handlerId, err := api.ParseUrlParamsId[int](ctx, "handlerId")
	if err != nil {
		return err
	}
	behaviorId, err := api.ParseUrlParamsId[int](ctx, "id")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.UpdateStubBehaviorForm](ctx)
	if err != nil {
		return err
	}
	err = a.useCase.UpdateRestStubBehavior(ctx.Context(), serviceId, handlerId, behaviorId, form)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// updateRestMockBehavior godoc
// @ID updateRestMockBehavior
// @Summary update existing mock behavior
// @Tags mockservices
// @Accept json
// @Produce json
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Param handlerId path number true "id of rest handler"
// @Param id path number true "id of behavior to update"
// @Route /mockservices/rest/{serviceId}/handlers/{handlerId}/updateRestMockBehavior/{id} [put]
func (a *API) updateRestMockBehavior(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	handlerId, err := api.ParseUrlParamsId[int](ctx, "handlerId")
	if err != nil {
		return err
	}
	behaviorId, err := api.ParseUrlParamsId[int](ctx, "id")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.UpdateMockBehaviorForm](ctx)
	if err != nil {
		return err
	}
	err = a.useCase.UpdateRestMockBehavior(ctx.Context(), serviceId, handlerId, behaviorId, form)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// moveRestBehaviorPriority godoc
// @ID moveRestBehaviorPriority
// @Summary move behavior in list to change priority of execution
// @Tags mockservices
// @Accept json
// @Produce json
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Param handlerId path number true "id of rest handler"
// @Param priority query number true "new priority of a rest handler behavior"
// @Route /mockservices/rest/{serviceId}/handlers/{handlerId}/behaviors/{id}/move [patch]
func (a *API) moveRestBehaviorPriority(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	handlerId, err := api.ParseUrlParamsId[int](ctx, "handlerId")
	if err != nil {
		return err
	}
	behaviorId, err := api.ParseUrlParamsId[int](ctx, "id")
	if err != nil {
		return err
	}
	newPriority, err := api.ParseQueryParam[int](ctx, "priority")
	if err != nil {
		return err
	}
	err = a.useCase.MoveRestBehaviorPriority(ctx.Context(), serviceId, handlerId, behaviorId, newPriority)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// deleteRestBehavior godoc
// @ID deleteRestBehavior
// @Summary delete behavior of rest route
// @Tags mockservices
// @Accept json
// @Produce json
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "id of a project where service is located"
// @Param serviceId path string true "id of a rest service"
// @Param handlerId path number true "id of rest handler"
// @Route /mockservices/rest/{serviceId}/handlers/{handlerId}/behaviors/{id} [delete]
func (a *API) deleteRestBehavior(ctx *fiber.Ctx) error {
	serviceId, err := parseServiceId(ctx)
	if err != nil {
		return err
	}
	handlerId, err := api.ParseUrlParamsId[int](ctx, "handlerId")
	if err != nil {
		return err
	}
	behaviorId, err := api.ParseUrlParamsId[int](ctx, "id")
	if err != nil {
		return err
	}
	err = a.useCase.DeleteRestBehavior(ctx.Context(), serviceId, handlerId, behaviorId)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

func parseServiceId(ctx *fiber.Ctx) (mockservice.RestServiceIdentifier, error) {
	serviceId, err := api.ParseUrlParamsId[string](ctx, "id")
	if err != nil {
		return mockservice.RestServiceIdentifier{}, err
	}
	projectId, err := api.ParseQueryParam[entity.BigIntPK](ctx, "projectId")
	if err != nil {
		return mockservice.RestServiceIdentifier{}, err
	}
	return mockservice.RestServiceIdentifier{
		ProjectId: projectId,
		ServiceId: serviceId,
	}, nil
}
