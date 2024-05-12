package mockservice

import (
	mockserviceRules "configurator/internal/businessrules/mockservice"
	"configurator/internal/entity"
	"configurator/pkg/exportstruct"
	"context"
	"errors"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

type redisPubSubUseCase struct {
	log      *zap.SugaredLogger
	projRepo IProjectRepository
}

type CreateRedisPubSubForm struct {
	Id   string `json:"id"`
	Port int    `json:"port"`
}

func (uc *redisPubSubUseCase) CreateRedisPubSub(ctx context.Context, projectId entity.BigIntPK, form CreateRedisPubSubForm) error {
	return uc.projRepo.ModifyProjectData(ctx, projectId, func(data *exportstruct.Config) error {
		if err := mockserviceRules.ValidateRestUniqServiceId(*data, form.Id); err != nil {
			return err
		}
		if err := mockserviceRules.ValidateUniqServicePort(*data, form.Port); err != nil {
			return err
		}

		data.Brokers = append(data.Brokers, exportstruct.Broker{
			ID:   form.Id,
			Type: exportstruct.BrokerType_REDIS_PUBSUB,
			Port: exportstruct.Port{
				Port: form.Port,
			},
		})
		return nil
	})
}

type UpdateRedisPubSubForm struct {
	Id   string `json:"id"`
	Port int    `json:"port"`
}

type BrokerIdentifier struct {
	ProjectId entity.BigIntPK
	BrokerId  string
}

func (uc *redisPubSubUseCase) UpdateRedisPubSub(ctx context.Context, id BrokerIdentifier, form UpdateRedisPubSubForm) error {
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		_, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}
		data.Brokers[brokerIdx].ID = form.Id
		data.Brokers[brokerIdx].Port = exportstruct.Port{
			Port: form.Port,
		}
		return nil
	})
}

type ListBrokerTopicsResultItem struct {
	Id    int    `json:"id"`
	Topic string `json:"topic"`
}

func (uc *redisPubSubUseCase) ListTopics(ctx context.Context, id BrokerIdentifier) ([]ListBrokerTopicsResultItem, error) {
	projData, err := fetchProjectData(ctx, uc.projRepo, id.ProjectId)
	if err != nil {
		return nil, err
	}
	broker, ok := lo.Find(projData.Brokers, func(item exportstruct.Broker) bool {
		return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
	})
	if !ok {
		return nil, errors.New("некорректный id брокера")
	}
	return lo.Map(broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior, func(behav exportstruct.BrokerBehaviorRedisItem, idx int) ListBrokerTopicsResultItem {
		return ListBrokerTopicsResultItem{
			Id:    idx,
			Topic: behav.Topic,
		}
	}), nil
}

type (
	ListTopicGeneratorsResult struct {
		Stubs []ListTopicGeneratorsStubItem `json:"stubs"`
		Mocks []ListTopicGeneratorsMockItem `json:"mocks"`
	}

	ListTopicGeneratorsStubItem struct {
		Id              int    `json:"id"`
		Priority        int    `json:"priority"`
		Interval        string `json:"interval"`
		SendImmediately bool   `json:"send_immediately"`

		Payload string `json:"payload"`
	}

	ListTopicGeneratorsMockItem struct {
		Id              int    `json:"id"`
		Priority        int    `json:"priority"`
		Interval        string `json:"interval"`
		SendImmediately bool   `json:"send_immediately"`

		Impl []string `json:"impl"`
	}
)

func (uc *redisPubSubUseCase) ListTopicGenerators(ctx context.Context, id BrokerIdentifier, topicId int) (ListTopicGeneratorsResult, error) {
	if topicId < 0 {
		return ListTopicGeneratorsResult{}, errors.New("некорректный идентификатор топика")
	}
	projData, err := fetchProjectData(ctx, uc.projRepo, id.ProjectId)
	if err != nil {
		return ListTopicGeneratorsResult{}, err
	}
	broker, ok := lo.Find(projData.Brokers, func(item exportstruct.Broker) bool {
		return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
	})
	if !ok {
		return ListTopicGeneratorsResult{}, errors.New("некорректный id брокера")
	}
	if topicId >= len(broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior) {
		return ListTopicGeneratorsResult{}, errors.New("некорректный идентификатор топика")
	}
	gens := broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior[topicId].Generators
	return ListTopicGeneratorsResult{
		Stubs: lo.FilterMap(gens, func(g exportstruct.RedisTopicGenerator, idx int) (ListTopicGeneratorsStubItem, bool) {
			return ListTopicGeneratorsStubItem{
				Id:              idx,
				Priority:        idx,
				Interval:        string(g.Interval),
				SendImmediately: g.SendImmediately,
				Payload:         g.RedisTopicGeneratorUnion.RedisTopicGeneratorConst.Payload,
			}, g.Type == exportstruct.RedisTopicGeneratorType_CONST
		}),
		Mocks: lo.FilterMap(gens, func(g exportstruct.RedisTopicGenerator, idx int) (ListTopicGeneratorsMockItem, bool) {
			return ListTopicGeneratorsMockItem{
				Id:              idx,
				Priority:        idx,
				Interval:        string(g.Interval),
				SendImmediately: g.SendImmediately,
				Impl:            g.RedisTopicGeneratorUnion.RedisTopicGeneratorProg.Behavior,
			}, g.Type == exportstruct.RedisTopicGeneratorType_PROG
		}),
	}, nil
}
