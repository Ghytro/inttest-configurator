package projects

import (
	"configurator/internal/entity"
	entProjects "configurator/internal/entity/projects"
	"configurator/internal/usecase/internal/common"
	"context"
)

type ProjectRepository interface {
	common.Repo

	ListProjects(ctx context.Context) (result []entProjects.Project, err error)
	CreateProject(ctx context.Context, project entProjects.Project) (resultEnt entProjects.Project, err error)
	GetProject(ctx context.Context, id entity.BigIntPK, fetchData bool) (result entProjects.Project, err error)
	UpdateProjectMeta(ctx context.Context, proj entProjects.Project) (result entProjects.Project, err error)
	DeleteProject(ctx context.Context, id entity.BigIntPK, softDelete bool) error
}
