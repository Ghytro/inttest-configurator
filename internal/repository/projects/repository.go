package projects

import (
	entProjects "configurator/internal/entity/projects"
	"configurator/internal/repository/internal/common"
	"configurator/pkg/database"
	"configurator/pkg/exportstruct"
	"context"
	"encoding/json"

	"github.com/go-pg/pg/v10/orm"
	"go.uber.org/zap"
)

type Repository struct {
	common.Mixin
}

func New(db *database.PGDB, log *zap.SugaredLogger) *Repository {
	return &Repository{
		Mixin: common.Mixin{
			DB:          db,
			Log:         log.With(zap.String("repo", "projects_repo")),
			ErrWrapDesc: "err in projects_repo",
		},
	}
}

func (r *Repository) ListProjects(ctx context.Context) (result []entProjects.Project, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		if err := db.Model(&result).Select(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) CreateProject(ctx context.Context, name string, desc *string, data exportstruct.Config) (resultEnt entProjects.Project, err error) {
	dataBin, err := json.Marshal(data)
	if err != nil {
		return entProjects.Project{}, err
	}
	resultEnt = entProjects.Project{
		Name:    name,
		Desc:    desc,
		RawData: dataBin,
	}
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		_, err := db.Model(&resultEnt).Returning("*").Insert()
		return err
	})
	if err != nil {
		return entProjects.Project{}, err
	}
	return resultEnt, nil
}
