package mockservice

import (
	"configurator/internal/entity"
	entProjects "configurator/internal/entity/projects"
	"configurator/pkg/exportstruct"
	"context"
)

type IProjectRepository interface {
	GetProject(ctx context.Context, id entity.BigIntPK, fetchData bool) (result entProjects.Project, err error)
	ModifyProjectData(ctx context.Context, projectId entity.BigIntPK, fn func(data *exportstruct.Config) error) error
}
