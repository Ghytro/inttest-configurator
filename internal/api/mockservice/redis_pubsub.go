package mockservice

import (
	"configurator/internal/api"
	"configurator/internal/entity"
	"configurator/internal/usecase/mockservice"

	"github.com/gofiber/fiber/v2"
)

func (a *API) registerRedisPubSubRoutes(router fiber.Router) {
	r := fiber.New()
	brokerIdGroup := r.Group(":/brokerId")
	brokerIdGroup.Post("/", a.createRedisPubSubBroker).
		Put("/", a.updateRedisPubSub).
		Delete("/", a.deleteRedisPubSub)

	topicsGroup := brokerIdGroup.Group("/topics").
		Get("/", a.listRedisPubSubTopics).
		Post("/", a.createRedisPubSubTopic)

	topicIdGroup := topicsGroup.Group("/:topicId").
		Put("/", a.updateRedisPubSubTopic).
		Delete("/", a.deleteRedisPubTopic).
		Post("/createConstGenerator/", a.createRedisPubSubConstTopicGenerator).
		Post("/createProgGenerator/", a.createRedisPubSubProgTopicGenerator).
		Put("/updConstGenerator/:id", a.updateRedisPubSubConstTopicGenerator).
		Put("updProgGenerator/:id", a.updateRedisPubSubProgTopicGenerator)

	topicIdGroup.Group("/generators").
		Get("/", a.listRedisPubSubTopicGenerators).
		Delete("/:id", a.deleteRedisPubSubGenerator)

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

// deleteRedisPubSub godoc
// @ID deleteRedisPubSub
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param brokerId path string true "id of broker to list topics"
// @Router /mockservices/redis-pubsub/{brokerId} [delete]
func (a *API) deleteRedisPubSub(ctx *fiber.Ctx) error {
	id, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	if err := a.useCase.DeleteRedisPubSub(ctx.Context(), id); err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// createRedisPubSubTopic godoc
// @ID createRedisPubSubTopic
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param brokerId path string true "id of broker to list topics"
// @Param form body mockservice.CreateTopicReq true "body form"
// @Router /mockservices/redis-pubsub/{brokerId}/topics [post]
func (a *API) createRedisPubSubTopic(ctx *fiber.Ctx) error {
	id, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.CreateTopicReq](ctx)
	if err != nil {
		return err
	}
	if err := a.useCase.CreateTopic(ctx.Context(), id, form); err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// updateRedisPubSubTopic godoc
// @ID updateRedisPubSubTopic
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param brokerId path string true "id of broker to list topics"
// @Param topicId path number true "id of topic to list generators"
// @Router /mockservices/redis-pubsub/{brokerId}/topics/{topicId} [put]
func (a *API) updateRedisPubSubTopic(ctx *fiber.Ctx) error {
	id, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.UpdTopicReq](ctx)
	if err != nil {
		return err
	}
	topicId, err := api.ParseUrlParamsId[int](ctx, "topicId")
	if err != nil {
		return err
	}
	if err := a.useCase.UpdateTopic(ctx.Context(), id, topicId, form); err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// deleteRedisPubTopic godoc
// @ID deleteRedisPubTopic
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param brokerId path string true "id of broker to list topics"
// @Param topicId path number true "id of topic to list generators"
// @Router /mockservices/redis-pubsub/{brokerId}/topics/{topicId} [delete]
func (a *API) deleteRedisPubTopic(ctx *fiber.Ctx) error {
	id, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	topicId, err := api.ParseUrlParamsId[int](ctx, "topicId")
	if err != nil {
		return err
	}
	if err := a.useCase.DeleteTopic(ctx.Context(), id, topicId); err != nil {
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

// createRedisPubSubConstTopicGenerator godoc
// @ID createRedisPubSubConstTopicGenerator
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param brokerId path string true "id of broker to list topics"
// @Param topicId path number true "id of topic to list generators"
// @Router /mockservices/redis-pubsub/{brokerId}/topics/{topicId}/createConstGenerator [post]
func (a *API) createRedisPubSubConstTopicGenerator(ctx *fiber.Ctx) error {
	brokerId, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	topicId, err := api.ParseUrlParamsId[int](ctx, "topicId")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.CreateStubTopicGeneratorReq](ctx)
	if err != nil {
		return err
	}
	if err := a.useCase.CreateConstTopicGenerator(ctx.Context(), brokerId, topicId, form); err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// createRedisPubSubProgTopicGenerator godoc
// @ID createRedisPubSubProgTopicGenerator
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param brokerId path string true "id of broker to list topics"
// @Param topicId path number true "id of topic to list generators"
// @Router /mockservices/redis-pubsub/{brokerId}/topics/{topicId}/createProgGenerator [post]
func (a *API) createRedisPubSubProgTopicGenerator(ctx *fiber.Ctx) error {
	brokerId, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	topicId, err := api.ParseUrlParamsId[int](ctx, "topicId")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.CreateMockTopicGeneratorReq](ctx)
	if err != nil {
		return err
	}
	if err := a.useCase.CreateProgTopicGenerator(ctx.Context(), brokerId, topicId, form); err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// updateRedisPubSubConstTopicGenerator godoc
// @ID updateRedisPubSubConstTopicGenerator
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param brokerId path string true "id of broker to list topics"
// @Param topicId path number true "id of topic to list generators"
// @Param id path number true "id of generator to delete"
// @Router /mockservices/redis-pubsub/{brokerId}/topics/{topicId}/updConstGenerator/{id} [put]
func (a *API) updateRedisPubSubConstTopicGenerator(ctx *fiber.Ctx) error {
	brokerId, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	topicId, err := api.ParseUrlParamsId[int](ctx, "topicId")
	if err != nil {
		return err
	}
	genId, err := api.ParseUrlParamsId[int](ctx, "id")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.UpdConstTopicGeneratorReq](ctx)
	if err != nil {
		return err
	}
	if err := a.useCase.UpdateConstTopicGenerator(ctx.Context(), brokerId, topicId, genId, form); err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// updateRedisPubSubProgTopicGenerator godoc
// @ID updateRedisPubSubProgTopicGenerator
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param brokerId path string true "id of broker to list topics"
// @Param topicId path number true "id of topic to list generators"
// @Param id path number true "id of generator to delete"
// @Router /mockservices/redis-pubsub/{brokerId}/topics/{topicId}/updProgGenerator/{id} [put]
func (a *API) updateRedisPubSubProgTopicGenerator(ctx *fiber.Ctx) error {
	brokerId, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	topicId, err := api.ParseUrlParamsId[int](ctx, "topicId")
	if err != nil {
		return err
	}
	genId, err := api.ParseUrlParamsId[int](ctx, "id")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[mockservice.UpdMockTopicGeneratorReq](ctx)
	if err != nil {
		return err
	}
	if err := a.useCase.UpdateProgTopicGenerator(ctx.Context(), brokerId, topicId, genId, form); err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// deleteRedisPubSubGenerator godoc
// @ID deleteRedisPubSubGenerator
// @Summary "update redis pubsub broker data"
// @Tags mockservices
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param projectId query number true "project id to create redis pubsub mock"
// @Param brokerId path string true "id of broker to list topics"
// @Param topicId path number true "topic id to delete generator from"
// @Param id path number true "id of generator to delete"
// @Router /mockservices/redis-pubsub/{brokerId}/topics/{topicId}/generators/{id} [delete]
func (a *API) deleteRedisPubSubGenerator(ctx *fiber.Ctx) error {
	brokerId, err := parseBrokerIdentifier(ctx)
	if err != nil {
		return err
	}
	topicId, err := api.ParseUrlParamsId[int](ctx, "topicId")
	if err != nil {
		return err
	}
	genId, err := api.ParseUrlParamsId[int](ctx, "id")
	if err != nil {
		return err
	}
	if err := a.useCase.DeleteTopicGenerator(ctx.Context(), brokerId, topicId, genId); err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

func parseBrokerIdentifier(ctx *fiber.Ctx) (mockservice.BrokerIdentifier, error) {
	projectId, err := api.ParseQueryParam[entity.BigIntPK](ctx, "projectId")
	if err != nil {
		return mockservice.BrokerIdentifier{}, err
	}
	brokerId, err := api.ParseUrlParamsId[string](ctx, "brokerId")
	if err != nil {
		return mockservice.BrokerIdentifier{}, err
	}
	return mockservice.BrokerIdentifier{
		ProjectId: projectId,
		BrokerId:  brokerId,
	}, nil
}
