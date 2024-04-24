package projects

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	projectRules "configurator/internal/businessrules/projects"
	"configurator/internal/entity"
	entProjects "configurator/internal/entity/projects"
	"configurator/pkg/exportstruct"

	"go.uber.org/zap"
)

type UseCase struct {
	log      *zap.SugaredLogger
	projRepo ProjectRepository
}

func New(projectRepo ProjectRepository, log *zap.SugaredLogger) *UseCase {
	return &UseCase{
		log:      log,
		projRepo: projectRepo,
	}
}

func (uc *UseCase) UploadProject(ctx context.Context, projectName string, projectDesc *string, uploadedConfig io.ReadCloser) (result entProjects.Project, err error) {
	err = uc.projRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := projectRules.ValidateProjectNameUniq(ctx, uc.projRepo, projectName); err != nil {
			return err
		}
		confRaw, err := io.ReadAll(uploadedConfig)
		if err != nil {
			return err
		}
		uploadedConfig.Close()
		var parsedConfig exportstruct.Config
		if err := json.Unmarshal(confRaw, &parsedConfig); err != nil {
			return err
		}
		uploadedConfig.Close()
		if err := parsedConfig.Validate(); err != nil {
			return err
		}
		createdProj := entProjects.Project{
			Name:    projectName,
			Desc:    projectDesc,
			RawData: confRaw,
		}
		result, err = uc.projRepo.CreateProject(ctx, createdProj)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entProjects.Project{}, err
	}
	return result, nil
}

func (uc *UseCase) UnloadProject(ctx context.Context, id entity.BigIntPK) (projectFilePath, projectName string, err error) {
	err = uc.projRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		proj, err := uc.projRepo.GetProject(ctx, id, true)
		if err != nil {
			return err
		}
		var conf exportstruct.Config
		if err := json.Unmarshal(proj.RawData, &conf); err != nil {
			return err
		}
		exportFile, err := os.CreateTemp("/tmp", "configurator-export-*")
		if err != nil {
			return err
		}
		defer exportFile.Close()
		if err := json.NewEncoder(exportFile).Encode(&conf); err != nil {
			return err
		}
		projectFilePath = exportFile.Name()
		projectName = proj.Name
		return nil
	})
	if err != nil {
		return "", "", err
	}
	return projectFilePath, projectName, nil
}

func (uc *UseCase) DeleteProject(ctx context.Context, id entity.BigIntPK) error {
	return uc.projRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		_, err := uc.projRepo.GetProject(ctx, id, false)
		if err != nil {
			return errors.New("некорректный идентификатор проекта")
		}
		return uc.projRepo.DeleteProject(ctx, id, false)
	})
}

func (uc *UseCase) UpdateProjectMeta(ctx context.Context, id entity.BigIntPK, newName string, newDesc *string) (result entProjects.Project, err error) {
	err = uc.projRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		proj, err := uc.projRepo.GetProject(ctx, id, false)
		if err != nil {
			return err
		}
		proj.Name = newName
		proj.Desc = newDesc
		result, err = uc.projRepo.UpdateProjectMeta(ctx, proj)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entProjects.Project{}, err
	}
	return result, nil
}

func (uc *UseCase) ListProjects(ctx context.Context) (result []entProjects.Project, err error) {
	err = uc.projRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		result, err = uc.projRepo.ListProjects(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
