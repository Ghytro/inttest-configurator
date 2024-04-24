package projects

import (
	"context"
	"io"

	"configurator/internal/entity"
	entProjects "configurator/internal/entity/projects"
)

type IUseCase interface {
	ListProjects(ctx context.Context) (result []entProjects.Project, err error)
	UpdateProjectMeta(ctx context.Context, id entity.BigIntPK, newName string, newDesc *string) (result entProjects.Project, err error)
	UploadProject(ctx context.Context, projectName string, projectDesc *string, uploadedConfig io.ReadCloser) (result entProjects.Project, err error)
	UnloadProject(ctx context.Context, id entity.BigIntPK) (projectFilePath, projectName string, err error)
	DeleteProject(ctx context.Context, id entity.BigIntPK) error
}
