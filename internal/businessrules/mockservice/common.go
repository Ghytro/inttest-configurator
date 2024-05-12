package mockservice

import (
	"configurator/pkg/exportstruct"
	"fmt"

	"github.com/samber/lo"
)

func ValidateUniqServicePort(projData exportstruct.Config, port int) error {
	if lo.ContainsBy(projData.RpcServices, func(item exportstruct.RpcService) bool {
		return item.GetPort() == port
	}) || lo.ContainsBy(projData.Brokers, func(item exportstruct.Broker) bool {
		return item.GetPort() == port
	}) {
		return fmt.Errorf("невозможно присвоить порт %d сервису: занято", port)
	}
	return nil
}
