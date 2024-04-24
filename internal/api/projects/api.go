package projects

import (
	"configurator/internal/api"
	"configurator/internal/entity"
	"configurator/internal/entity/auth"
	"errors"
	"io"
	"os"

	entProjects "configurator/internal/entity/projects"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type API struct {
	log     *zap.SugaredLogger
	useCase IUseCase
}

func New(useCase IUseCase, log *zap.SugaredLogger) *API {
	return &API{
		log:     log,
		useCase: useCase,
	}
}

func (a *API) Register(router fiber.Router, authenticator api.Authenticator, middlewares ...fiber.Handler) {
	r := fiber.New()
	r.Use("/", authenticator.AuthMiddleware())
	for _, m := range middlewares {
		r.Use("/", m)
	}
	r.Get(
		"/",
		authenticator.CheckPermsMiddleware(auth.PermissionCreateProject),
		a.listProjects,
	)
	r.Post(
		"/",
		authenticator.CheckPermsMiddleware(auth.PermissionCreateProject),
		a.createProject,
	)
	r.Post(
		"/:id/unload",
		authenticator.CheckPermsMiddleware(auth.PermissionEditProject),
		a.unloadProject,
	)
	r.Put(
		"/:id",
		authenticator.CheckPermsMiddleware(auth.PermissionEditProject),
		a.updateProjectMeta,
	)
	r.Delete(
		"/:id",
		authenticator.CheckPermsMiddleware(auth.PermissionEditProject),
		a.deleteProject,
	)

	router.Mount("/projects", r)
}

// listProjects godoc
// @ID listProjects
// @Summary list all the projects in configurator
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {array} listProjectResponseItem
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Router /projects [get]
func (a *API) listProjects(ctx *fiber.Ctx) error {
	projects, err := a.useCase.ListProjects(ctx.Context())
	if err != nil {
		return err
	}
	response := lo.Map(projects, func(p entProjects.Project, _ int) listProjectResponseItem {
		return listProjectResponseItem{
			Id:   p.Id,
			Name: p.Name,
			Desc: p.Desc,
		}
	})
	return ctx.JSON(response)
}

// @createProject godoc
// @ID createProject
// @Summary create a new project
// @Tags projects
// @Accept mpfd
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param file formData file false "file with config to create project based on it"
// @Param name formData string true "name of created project"
// @Param desc formData string false "description of a created project"
// @Router /projects [post]
func (a *API) createProject(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}
	var (
		projFileReader io.ReadCloser
		projName       string
		projDesc       *string
	)
	{
		if projectFiles, ok := form.File["file"]; ok && len(projectFiles) != 0 {
			projFileReader, err = projectFiles[0].Open()
			if err != nil {
				return err
			}
		}
		projNames, ok := form.Value["name"]
		if !ok || len(projNames) == 0 {
			return errors.New("имя проекта обязательно")
		}
		projName = projNames[0]
		if projDescs, ok := form.Value["desc"]; !ok && len(projDescs) != 0 {
			projDesc = lo.ToPtr(projDescs[0])
		}
	}

	_, err = a.useCase.UploadProject(ctx.Context(), projName, projDesc, projFileReader)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// updateProjectMeta godoc
// @ID updateProjectMeta
// @Summary update name and desc of a project
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param id path number true "id of an updated project"
// @Param form body updateProjectMetaRequest true "form to update project meta"
// @Param id path number true "id of updated project"
// @Router /projects/{id} [put]
func (a *API) updateProjectMeta(ctx *fiber.Ctx) error {
	id, err := api.ParseUrlParamsId[entity.BigIntPK](ctx, "id")
	if err != nil {
		return err
	}
	form, err := api.ParseBody[updateProjectMetaReqest](ctx)
	if err != nil {
		return err
	}
	_, err = a.useCase.UpdateProjectMeta(ctx.Context(), id, form.Name, form.Desc)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// unloadProject godoc
// @ID unloadProject
// @Summary form a json config, unload project from configurator
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {string} string "configuration file"
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param id path number true "id of a project to unload"
// @Router /projects/{id}/unload [post]
func (a *API) unloadProject(ctx *fiber.Ctx) error {
	projectId, err := api.ParseUrlParamsId[entity.BigIntPK](ctx, "id")
	if err != nil {
		return err
	}
	filePath, projectName, err := a.useCase.UnloadProject(ctx.Context(), projectId)
	if err != nil {
		return err
	}
	defer os.Remove(filePath)
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return ctx.Download(filePath, projectName+".json")
}

// deleteProject godoc
// @ID deleteProject
// @Summary delete a project forever
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param id path number true "id of a project to delete"
// @Router /projects/{id} [delete]
func (a *API) deleteProject(ctx *fiber.Ctx) error {
	id, err := api.ParseUrlParamsId[entity.BigIntPK](ctx, "id")
	if err != nil {
		return err
	}
	err = a.useCase.DeleteProject(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}
