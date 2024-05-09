package projects

import (
	"configurator/internal/entity"
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

func (r *Repository) CreateProject(ctx context.Context, project entProjects.Project) (resultEnt entProjects.Project, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		resultEnt = project
		_, err := db.Model(&resultEnt).Returning("*").Insert()
		return err
	})
	if err != nil {
		return entProjects.Project{}, err
	}
	return resultEnt, nil
}

func (r *Repository) GetProject(ctx context.Context, id entity.BigIntPK, fetchData bool) (result entProjects.Project, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		columns := []string{"id", "name", "description"}
		if fetchData {
			columns = append(columns, "data")
		}
		result.Id = id
		if err := db.Model(&result).Column(columns...).WherePK().Select(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entProjects.Project{}, err
	}
	return result, nil
}

func (r *Repository) UpdateProjectMeta(ctx context.Context, proj entProjects.Project) (result entProjects.Project, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		result = proj
		_, err := db.Model(&result).WherePK().
			Set("name = ?", result.Name).
			Set("desc = ?", result.Desc).
			Update()
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

func (r *Repository) DeleteProject(ctx context.Context, id entity.BigIntPK, softDelete bool) error {
	return r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		model := entProjects.Project{Id: id}
		if _, err := db.Model(&model).WherePK().Delete(); err != nil {
			return err
		}
		return nil
	})
}

func (r *Repository) UpdateProjectData(ctx context.Context, projectId entity.BigIntPK, newData exportstruct.Config) error {
	rawData, err := json.Marshal(newData)
	if err != nil {
		return err
	}
	return r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		project, err := r.GetProject(ctx, projectId, false)
		if err != nil {
			return err
		}
		project.RawData = rawData
		if _, err := db.Model(&project).WherePK().Update(); err != nil {
			return err
		}
		return nil
	})
}

func (r *Repository) ModifyProjectData(ctx context.Context, projectId entity.BigIntPK, fn func(data *exportstruct.Config) error) error {
	proj, err := r.GetProject(ctx, projectId, true)
	if err != nil {
		return err
	}
	data, err := proj.ParsedData()
	if err != nil {
		return err
	}
	if err := fn(&data); err != nil {
		return err
	}
	if err := r.UpdateProjectData(ctx, projectId, data); err != nil {
		return err
	}
	return nil
}
