package projects

import (
	entProjects "configurator/internal/entity/projects"
	"context"
	"errors"

	"github.com/samber/lo"
)

type ProjectsReader interface {
	ListProjects(ctx context.Context) ([]entProjects.Project, error)
}

func ValidateProjectNameUniq(ctx context.Context, reader ProjectsReader, createdName string) error {
	allProj, err := reader.ListProjects(ctx)
	if err != nil {
		return err
	}
	if lo.ContainsBy(allProj, func(item entProjects.Project) bool {
		return item.Name == createdName
	}) {
		return errors.New("наименование создаваемого проекта должно быть уникально")
	}
	return nil
}
