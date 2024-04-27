package mockservice

import (
	"configurator/internal/entity"
	"configurator/pkg/exportstruct"
	"context"

	"go.uber.org/zap"
)

type UseCase struct {
	log      *zap.SugaredLogger
	projRepo IProjectRepository
	restServiceUseCase
}

func New(log *zap.SugaredLogger, projRepo IProjectRepository) *UseCase {
	return &UseCase{
		log:      log,
		projRepo: projRepo,
		restServiceUseCase: restServiceUseCase{
			log:      log,
			projRepo: projRepo,
		},
	}
}

func (uc *UseCase) ListServices(ctx context.Context, projId entity.BigIntPK) ([]CommonServiceResult, error) {
	data, err := fetchProjectData(ctx, uc.projRepo, projId)
	if err != nil {
		return nil, err
	}

	result := make([]CommonServiceResult, 0, len(data.RpcServices)+len(data.Brokers))
	for _, s := range data.RpcServices {
		result = append(result, CommonServiceResult{
			Id:   s.ID,
			Port: s.GetPort(),
			Type: string(s.Type),
		})
	}
	for _, b := range data.Brokers {
		result = append(result, CommonServiceResult{
			Id:   b.ID,
			Port: b.GetPort(),
			Type: string(b.Type),
		})
	}
	return result, nil
}

func fetchProjectData(ctx context.Context, repo IProjectRepository, projId entity.BigIntPK) (exportstruct.Config, error) {
	project, err := repo.GetProject(ctx, projId, true)
	if err != nil {
		return exportstruct.Config{}, err
	}
	data, err := project.ParsedData()
	if err != nil {
		return exportstruct.Config{}, err
	}
	return data, nil
}
