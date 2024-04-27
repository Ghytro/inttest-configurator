package mockservice

import (
	"configurator/pkg/exportstruct"
	"fmt"

	"github.com/samber/lo"
)

func ValidateRestUniqServiceId(projData exportstruct.Config, createdId string) error {
	if lo.ContainsBy(projData.RpcServices, func(item exportstruct.RpcService) bool {
		return item.Type == exportstruct.RpcServiceType_REST && item.ID == createdId
	}) {
		return fmt.Errorf("уже существует REST-сервис с id=%s", createdId)
	}
	return nil
}
