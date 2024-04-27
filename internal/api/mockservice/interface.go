package mockservice

import (
	"configurator/internal/entity"
	"configurator/internal/usecase/mockservice"
	"context"
)

type IUseCase interface {
	ListServices(ctx context.Context, projId entity.BigIntPK) ([]mockservice.CommonServiceResult, error)
}
