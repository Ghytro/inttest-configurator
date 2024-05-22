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

func (uc *redisPubSubUseCase) DeleteRedisPubSub(ctx context.Context, id BrokerIdentifier) error {
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		_, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}
		data.Brokers = append(data.Brokers[brokerIdx:], data.Brokers[brokerIdx+1:]...)
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

type CreateTopicReq struct {
	Topic string `json:"topic"`
}

func (uc *redisPubSubUseCase) CreateTopic(ctx context.Context, id BrokerIdentifier, form CreateTopicReq) error {
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		broker, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}

		if err := mockserviceRules.ValidateUniqTopicInBroker(broker, form.Topic); err != nil {
			return err
		}

		brokerPtr := &data.Brokers[brokerIdx]
		behavs := broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior
		brokerPtr.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior = append(
			behavs,
			exportstruct.BrokerBehaviorRedisItem{
				Topic: form.Topic,
			},
		)
		return nil
	})
}

type UpdTopicReq struct {
	Name string `json:"name"`
}

func (uc *redisPubSubUseCase) UpdateTopic(ctx context.Context, id BrokerIdentifier, topicId int, form UpdTopicReq) error {
	if topicId < 0 {
		return errors.New("некорректный id топика")
	}
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		_, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}
		brokerPtr := &data.Brokers[brokerIdx]
		if topicId >= len(brokerPtr.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior) {
			return errors.New("некорректный идентификатор топика")
		}
		behaviorPtr := &brokerPtr.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior[topicId]
		behaviorPtr.Topic = form.Name
		return nil
	})
}

func (uc *redisPubSubUseCase) DeleteTopic(ctx context.Context, id BrokerIdentifier, topicId int) error {
	if topicId < 0 {
		return errors.New("некорректный id топика")
	}
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		broker, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}
		if topicId >= len(broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior) {
			return errors.New("некорректный идентификатор топика")
		}
		behavs := broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior
		behavs = append(behavs[:topicId], behavs[topicId+1:]...)
		data.Brokers[brokerIdx].BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior = behavs
		return nil
	})
}

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

type CreateStubTopicGeneratorReq struct {
	Interval        exportstruct.StrTimeInterval `json:"interval"`
	Payload         string                       `json:"payload"`
	SendImmediately bool                         `json:"send_immediately"`
}

func (uc *redisPubSubUseCase) CreateConstTopicGenerator(ctx context.Context, id BrokerIdentifier, topicId int, form CreateStubTopicGeneratorReq) error {
	if topicId < 0 {
		return errors.New("некорректный идентификатор топика")
	}
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		broker, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}
		if topicId >= len(broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior) {
			return errors.New("некорректный идентификатор топика")
		}

		brokerPtr := &data.Brokers[brokerIdx]
		behaviorPtr := brokerPtr.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior[topicId]
		behaviorPtr.Generators = append(behaviorPtr.Generators, exportstruct.RedisTopicGenerator{
			Interval:        form.Interval,
			SendImmediately: form.SendImmediately,
			Type:            exportstruct.RedisTopicGeneratorType_CONST,
			RedisTopicGeneratorUnion: exportstruct.RedisTopicGeneratorUnion{
				RedisTopicGeneratorConst: exportstruct.RedisTopicGeneratorConst{
					Payload: form.Payload,
				},
			},
		})
		return nil
	})
}

type UpdConstTopicGeneratorReq struct {
	Interval        exportstruct.StrTimeInterval `json:"interval"`
	Payload         string                       `json:"payload"`
	SendImmediately bool                         `json:"send_immediately"`
}

func (uc *redisPubSubUseCase) UpdateConstTopicGenerator(ctx context.Context, id BrokerIdentifier, topicId int, genId int, form UpdConstTopicGeneratorReq) error {
	if topicId < 0 {
		return errors.New("некорректный id топика")
	}
	if genId < 0 {
		return errors.New("некорректный id генератора поведения топика")
	}
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		_, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}
		brokerPtr := &data.Brokers[brokerIdx]

		if topicId >= len(brokerPtr.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior) {
			return errors.New("некорректный идентификатор топика")
		}
		behaviorPtr := brokerPtr.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior[topicId]

		if genId >= len(behaviorPtr.Generators) {
			return errors.New("некорректный id генератора поведения топика")
		}
		if behaviorPtr.Generators[genId].Type != exportstruct.RedisTopicGeneratorType_CONST {
			return errors.New("обновляемый генератор поведения топика не является константным")
		}

		behaviorPtr.Generators[genId] = exportstruct.RedisTopicGenerator{
			Interval:        form.Interval,
			SendImmediately: form.SendImmediately,
			Type:            exportstruct.RedisTopicGeneratorType_CONST,
			RedisTopicGeneratorUnion: exportstruct.RedisTopicGeneratorUnion{
				RedisTopicGeneratorConst: exportstruct.RedisTopicGeneratorConst{
					Payload: form.Payload,
				},
			},
		}
		return nil
	})
}

type CreateMockTopicGeneratorReq struct {
	Interval        exportstruct.StrTimeInterval `json:"interval"`
	Impl            []string                     `json:"impl"`
	SendImmediately bool                         `json:"send_immediately"`
}

func (uc *redisPubSubUseCase) CreateProgTopicGenerator(ctx context.Context, id BrokerIdentifier, topicId int, form CreateMockTopicGeneratorReq) error {
	if topicId < 0 {
		return errors.New("некорректный идентификатор топика")
	}
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		broker, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}
		if topicId >= len(broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior) {
			return errors.New("некорректный идентификатор топика")
		}

		brokerPtr := &data.Brokers[brokerIdx]
		behaviorPtr := brokerPtr.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior[topicId]
		behaviorPtr.Generators = append(behaviorPtr.Generators, exportstruct.RedisTopicGenerator{
			Interval:        form.Interval,
			SendImmediately: form.SendImmediately,
			Type:            exportstruct.RedisTopicGeneratorType_PROG,
			RedisTopicGeneratorUnion: exportstruct.RedisTopicGeneratorUnion{
				RedisTopicGeneratorProg: exportstruct.RedisTopicGeneratorProg{
					Behavior: form.Impl,
				},
			},
		})
		return nil
	})
}

type UpdMockTopicGeneratorReq struct {
	Interval        exportstruct.StrTimeInterval `json:"interval"`
	Impl            []string                     `json:"impl"`
	SendImmediately bool                         `json:"send_immediately"`
}

func (uc *redisPubSubUseCase) UpdateProgTopicGenerator(ctx context.Context, id BrokerIdentifier, topicId int, genId int, form UpdMockTopicGeneratorReq) error {
	if topicId < 0 {
		return errors.New("некорректный идентификатор топика")
	}
	if genId < 0 {
		return errors.New("некорректный id генератора поведения топика")
	}
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		broker, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}
		if topicId >= len(broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior) {
			return errors.New("некорректный идентификатор топика")
		}

		brokerPtr := &data.Brokers[brokerIdx]
		behaviorPtr := brokerPtr.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior[topicId]
		if genId >= len(behaviorPtr.Generators) {
			return errors.New("некорректный id генератора поведения топика")
		}
		if behaviorPtr.Generators[genId].Type != exportstruct.RedisTopicGeneratorType_PROG {
			return errors.New("обновляемый генератор поведения топика не программируемый")
		}
		behaviorPtr.Generators[genId] = exportstruct.RedisTopicGenerator{
			Interval:        form.Interval,
			SendImmediately: form.SendImmediately,
			Type:            exportstruct.RedisTopicGeneratorType_PROG,
			RedisTopicGeneratorUnion: exportstruct.RedisTopicGeneratorUnion{
				RedisTopicGeneratorProg: exportstruct.RedisTopicGeneratorProg{
					Behavior: form.Impl,
				},
			},
		}
		return nil
	})
}

func (uc *redisPubSubUseCase) DeleteTopicGenerator(ctx context.Context, id BrokerIdentifier, topicId int, genId int) error {
	if topicId < 0 {
		return errors.New("некорректный id топика")
	}
	if genId < 0 {
		return errors.New("некорректный id правила генератора топика")
	}
	return uc.projRepo.ModifyProjectData(ctx, id.ProjectId, func(data *exportstruct.Config) error {
		broker, brokerIdx, ok := lo.FindIndexOf(data.Brokers, func(item exportstruct.Broker) bool {
			return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == id.BrokerId
		})
		if !ok {
			return errors.New("некорректный id брокера")
		}
		brokerPtr := &data.Brokers[brokerIdx]
		if topicId >= len(broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior) {
			return errors.New("некорректный идентификатор топика")
		}
		behaviorPtr := brokerPtr.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior[topicId]
		if genId >= len(behaviorPtr.Generators) {
			return errors.New("некорректный id правила генератора топика")
		}
		behaviorPtr.Generators = append(behaviorPtr.Generators[:genId], behaviorPtr.Generators[genId+1:]...)
		return nil
	})
}
