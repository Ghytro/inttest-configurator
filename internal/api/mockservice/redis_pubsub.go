package mockservice

import (
	"configurator/internal/api"
	"configurator/internal/entity"
	"configurator/internal/usecase/mockservice"

	"github.com/gofiber/fiber/v2"
)

func (a *API) registerRedisPubSubRoutes(router fiber.Router) {
	r := fiber.New()
	r.Post("/", a.createRedisPubSubBroker).
		Put("/:id", a.updateRedisPubSub)

	topicsGroup := r.Group("/:brokerId/topics").
		Get("/", a.listRedisPubSubTopics)

	topicsGroup.Group("/:topicId").
		Get("/generators", a.listRedisPubSubTopicGenerators)

	router.Mount("/redis-pubsub", r)
}

// createRedisPubSubBroker godoc
// @ID createRedisPubSubBroker
// @Summary "create redis pubsub mock broker"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param form body mockservice.CreateRedisPubSubForm true "body form"
// @Router /mockservices/redis-pubsub [post]
func (a *API) createRedisPubSubBroker(ctx *fiber.Ctx) error {
	projectId, err := api.ParseQueryParam[entity.BigIntPK](ctx, "projectId")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.CreateRedisPubSubForm](ctx)
	if err != nil {
		return err
	}
	err = a.useCase.CreateRedisPubSub(ctx.Context(), projectId, form)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// listRedisPubSubTopics godoc
// @ID listRedisPubSubTopics
// @Summary "list redis pubsub broker topics"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {array} mockservice.ListBrokerTopicsResultItem
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to list all services"
// @Param brokerId path string true "id of broker to list topics"
// @Router /mockservices/redis-pubsub/{brokerId}/topics [get]
func (a *API) listRedisPubSubTopics(ctx *fiber.Ctx) error {
	id, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}

	result, err := a.useCase.ListTopics(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

// updateRedisPubSub godoc
// @ID updateRedisPubSub
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to list all services"
// @Param brokerId path string true "id of broker to update"
// @Param form body mockservice.UpdateRedisPubSubForm true "body form"
// @Router /mockservices/redis-pubsub/{brokerId} [put]
func (a *API) updateRedisPubSub(ctx *fiber.Ctx) error {
	id, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.UpdateRedisPubSubForm](ctx)
	if err != nil {
		return err
	}
	err = a.useCase.UpdateRedisPubSub(ctx.Context(), id, form)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// listRedisPubSubTopicGenerators godoc
// @ID listRedisPubSubTopicGenerators
// @Summary "list redis pubsub topic generators"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} mockservice.ListTopicGeneratorsResult
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to list all services"
// @Param brokerId path string true "id of broker to update"
// @Param topicId path number true "id of topic to list generators"
// @Router /mockservices/redis-pubsub/{brokerId}/topics/{topicId}/generators [get]
func (a *API) listRedisPubSubTopicGenerators(ctx *fiber.Ctx) error {
	brokerId, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	topicId, err := api.ParseUrlParamsId[int](ctx, "topicId")
	if err != nil {
		return err
	}
	result, err := a.useCase.ListTopicGenerators(ctx.Context(), brokerId, topicId)
	if err != nil {
		return err
	}
	return ctx.JSON(result)
}

func parseBrokerIdentifier(ctx *fiber.Ctx) (mockservice.BrokerIdentifier, error) {
	projectId, err := api.ParseQueryParam[entity.BigIntPK](ctx, "projectId")
	if err != nil {
		return mockservice.BrokerIdentifier{}, err
	}
	brokerId, err := api.ParseQueryParam[string](ctx, "brokerId")
	if err != nil {
		return mockservice.BrokerIdentifier{}, err
	}
	return mockservice.BrokerIdentifier{
		ProjectId: projectId,
		BrokerId:  brokerId,
	}, nil
}
