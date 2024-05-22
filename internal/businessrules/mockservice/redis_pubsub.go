package mockservice

import (
	"configurator/pkg/exportstruct"
	"fmt"

	"github.com/samber/lo"
)

func ValidateRedisPubSubUniqServiceId(projData exportstruct.Config, createdId string) error {
	if lo.ContainsBy(projData.Brokers, func(item exportstruct.Broker) bool {
		return item.Type == exportstruct.BrokerType_REDIS_PUBSUB && item.ID == createdId
	}) {
		return fmt.Errorf("уже существует REST-сервис с id=%s", createdId)
	}
	return nil
}

func ValidateUniqTopicInBroker(broker exportstruct.Broker, topic string) error {
	if lo.ContainsBy(
		broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior,
		func(item exportstruct.BrokerBehaviorRedisItem) bool {
			return item.Topic == topic
		},
	) {
		return fmt.Errorf("уже существует топик %s", topic)
	}
	return nil
}
