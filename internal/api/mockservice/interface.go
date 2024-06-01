package mockservice

import (
	"configurator/internal/entity"
	"configurator/internal/usecase/mockservice"
	"context"

	"github.com/pb33f/libopenapi"
)

type IUseCase interface {
	ListServices(ctx context.Context, projId entity.BigIntPK) ([]mockservice.CommonServiceResult, error)

	iRestServiceUseCase
	iRedisPubSubUsecase
}

type iRestServiceUseCase interface {
	CreateRestService(ctx context.Context, projectId entity.BigIntPK, request mockservice.CreateRestServiceReq) (result mockservice.CreateRestServiceResult, err error)
	UpdateRestService(ctx context.Context, serviceId mockservice.RestServiceIdentifier, request mockservice.RestServiceUpdateForm) (mockservice.RestServiceUpdateResult, error)
	DeleteRestService(ctx context.Context, serviceId mockservice.RestServiceIdentifier) error
	ListRestHandlers(ctx context.Context, serviceId mockservice.RestServiceIdentifier) ([]mockservice.ListRestRoutesResult, error)
	AddRestHandler(ctx context.Context, serviceId mockservice.RestServiceIdentifier, form mockservice.CreateRestHandlerForm) error
	UpdateRestHandler(ctx context.Context, serviceId mockservice.RestServiceIdentifier, id int, form mockservice.UpdateRestHandlerForm) error
	DeleteRestHandler(ctx context.Context, serviceId mockservice.RestServiceIdentifier, id int) error
	AddRestStubBehavior(ctx context.Context, serviceId mockservice.RestServiceIdentifier, handlerId int, form mockservice.CreateRestStubBehaviorForm) error
	AddRestMockBehavior(ctx context.Context, serviceId mockservice.RestServiceIdentifier, handlerId int, form mockservice.CreateRestMockBehaviorForm) error
	ListRestBehaviors(ctx context.Context, serviceId mockservice.RestServiceIdentifier, handlerId int) (result mockservice.ListRestBehaviorResult, err error)
	UpdateRestStubBehavior(ctx context.Context, serviceId mockservice.RestServiceIdentifier, handlerId, behaviorId int, form mockservice.UpdateStubBehaviorForm) error
	UpdateRestMockBehavior(ctx context.Context, serviceId mockservice.RestServiceIdentifier, handlerId, behaviorId int, form mockservice.UpdateMockBehaviorForm) error
	MoveRestBehaviorPriority(ctx context.Context, serviceId mockservice.RestServiceIdentifier, handlerId, behaviorId int, newPriority int) error
	DeleteRestBehavior(ctx context.Context, serviceId mockservice.RestServiceIdentifier, handlerId, behaviorId int) error
	ImportSwagger(ctx context.Context, serviceId mockservice.RestServiceIdentifier, swagger libopenapi.Document) error
}

type iRedisPubSubUsecase interface {
	CreateConstTopicGenerator(ctx context.Context, id mockservice.BrokerIdentifier, topicId int, form mockservice.CreateStubTopicGeneratorReq) error
	CreateProgTopicGenerator(ctx context.Context, id mockservice.BrokerIdentifier, topicId int, form mockservice.CreateMockTopicGeneratorReq) error
	CreateRedisPubSub(ctx context.Context, projectId entity.BigIntPK, form mockservice.CreateRedisPubSubForm) error
	CreateTopic(ctx context.Context, id mockservice.BrokerIdentifier, form mockservice.CreateTopicReq) error
	DeleteRedisPubSub(ctx context.Context, id mockservice.BrokerIdentifier) error
	DeleteTopic(ctx context.Context, id mockservice.BrokerIdentifier, topicId int) error
	DeleteTopicGenerator(ctx context.Context, id mockservice.BrokerIdentifier, topicId int, genId int) error
	ListTopicGenerators(ctx context.Context, id mockservice.BrokerIdentifier, topicId int) (mockservice.ListTopicGeneratorsResult, error)
	ListTopics(ctx context.Context, id mockservice.BrokerIdentifier) ([]mockservice.ListBrokerTopicsResultItem, error)
	UpdateConstTopicGenerator(ctx context.Context, id mockservice.BrokerIdentifier, topicId int, genId int, form mockservice.UpdConstTopicGeneratorReq) error
	UpdateProgTopicGenerator(ctx context.Context, id mockservice.BrokerIdentifier, topicId int, genId int, form mockservice.UpdMockTopicGeneratorReq) error
	UpdateRedisPubSub(ctx context.Context, id mockservice.BrokerIdentifier, form mockservice.UpdateRedisPubSubForm) error
	UpdateTopic(ctx context.Context, id mockservice.BrokerIdentifier, topicId int, form mockservice.UpdTopicReq) error
}
