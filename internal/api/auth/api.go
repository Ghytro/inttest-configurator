package auth

import (
	"errors"
	"time"

	"configurator/internal/api"
	"configurator/internal/entity"
	entAuth "configurator/internal/entity/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type API struct {
	log     *zap.SugaredLogger
	useCase IUseCase
}

func New(log *zap.SugaredLogger, useCase IUseCase) *API {
	return &API{
		log:     log.With(zap.String("api", "auth_api")),
		useCase: useCase,
	}
}

func (a *API) Register(router fiber.Router, authMiddleware fiber.Handler, middlewares ...fiber.Handler) {
	r := fiber.New()
	r.Post("/", a.auth)
	router.Mount("/auth", r)

	r = fiber.New()
	r.Use("/", authMiddleware)
	for _, m := range middlewares {
		r.Use(m)
	}
	r.Get("/", a.listUsers)
	r.Put("/:id", a.updateUser)
	r.Post("/", a.createUser)
	r.Delete("/:id", a.deleteUser)
	router.Mount("/users", r)

	r = fiber.New()
	r.Use("/", authMiddleware)
	for _, m := range middlewares {
		r.Use(m)
	}
	r.Get("/", a.listRoles)
	r.Post("/", a.createRole)
	r.Put("/:id", a.updateRole)
	r.Delete("/:id", a.deleteRole)
	router.Mount("/roles", r)

	r = fiber.New()
	r.Use("/", authMiddleware)
	r.Get("/", a.listPerms)
	router.Mount("/perms", r)
}

// createUser godoc
// @Summary create user
// @ID createUser
// @Tags users
// @Accept  	json
// @Produce  	json
// @Success 200 {object} createUserResponse
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param form body createUserRequest true "create user request model"
// @Router /users [post]
func (a *API) createUser(ctx *fiber.Ctx) error {
	const location = "ошибка создания пользователя"

	form, err := api.ParseBody[createUserRequest](ctx)
	if err != nil {
		return err
	}
	creator, ok := api.GetUserEntity(ctx)
	if !ok {
		return errors.New("user not authorized")
	}
	result, err := a.useCase.CreateUser(ctx.Context(), creator.Id, form.UserName, form.Password, form.Roles...)
	if err != nil {
		return err
	}
	return ctx.JSON(createUserResponse{
		Id:        result.Id,
		CreatedAt: result.CreatedAt,
	})
}

// updateUser godoc
// @ID updateUser
// @Summary update user
// @Tags users
// @Accept  	json
// @Produce  	json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param form body updateUserRequest true "create user request model"
// @Param id path number true "id of a user to update"
// @Router /users/{id} [put]
func (a *API) updateUser(ctx *fiber.Ctx) error {
	updater, ok := api.GetUserEntity(ctx)
	if !ok {
		return errors.New("user not authorized")
	}
	var model updateUserRequest
	if err := ctx.BodyParser(&model); err != nil {
		return err
	}
	userId, err := entity.ParseBigIntPK(ctx.Params("id"))
	if err != nil {
		return err
	}
	_, err = a.useCase.UpdateUser(ctx.Context(), updater.Id, userId, model.UserName, model.Roles...)
	if err != nil {
		return err
	}
	return ctx.JSON(api.OK{})
}

// deleteUser godoc
// @Summary delete user
// @ID deleteUser
// @Tags users
// @Accept  	json
// @Produce  	json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param id path number true "id of a user to delete"
// @Router /users/{id} [delete]
func (a *API) deleteUser(ctx *fiber.Ctx) error {
	const location = "ошибка удаления пользователя"

	id, err := api.ParseUrlParamsId[entity.BigIntPK](ctx, "id")
	if err != nil {
		return err
	}
	deleter, ok := api.GetUserEntity(ctx)
	if !ok {
		return errors.New("user not authorized")
	}
	err = a.useCase.DeleteUser(ctx.Context(), deleter.Id, id)
	if err != nil {
		return err
	}
	return ctx.JSON(nil)
}

// listUsers godoc
// @Summary create user
// @ID listUsers
// @Tags users
// @Accept  	json
// @Produce  	json
// @Success 200 {array} listUsersResponseItem
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Router /users [get]
func (a *API) listUsers(ctx *fiber.Ctx) error {
	const location = "ошибка получения списка пользователей"

	querier, ok := api.GetUserEntity(ctx)
	if !ok {
		return errors.New("no user entity in authorized ctx")
	}
	users, err := a.useCase.ListUsers(ctx.Context(), querier.Id)
	if err != nil {
		return err
	}
	return ctx.JSON(lo.Map(users, func(u entAuth.User, _ int) listUsersResponseItem {
		return listUsersResponseItem{
			Id:        u.Id,
			UserName:  u.UserName,
			CreatedAt: u.CreatedAt,
			Roles: lo.Map(u.Roles, func(r entAuth.Role, _ int) entity.BigIntPK {
				return r.Id
			}),
		}
	}))
}

// auth godoc
// @Summary create user
// @ID auth
// @Tags users
// @Accept  	json
// @Produce  	json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Param form body authRequest true "login/pass form"
// @Router /auth [post]
func (a *API) auth(ctx *fiber.Ctx) error {
	form, err := api.ParseBody[authRequest](ctx)
	if err != nil {
		return err
	}
	token, err := a.useCase.GenToken(ctx.Context(), form.UserName, form.Password)
	if err != nil {
		return err
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     api.JwtCookieName,
		Value:    token,
		Path:     api.JwtCookiePath,
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	return ctx.Send(nil)
}

// createRole godoc
// @ID createRole
// @Summary create role
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} roleCreateResponse
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param form body roleCreateRequest true "create role form"
// @Router /roles [post]
func (a *API) createRole(ctx *fiber.Ctx) error {
	user, ok := api.GetUserEntity(ctx)
	if !ok {
		return errors.New("user unauthorized")
	}
	var model roleCreateRequest
	if err := ctx.BodyParser(&model); err != nil {
		return err
	}
	createdRole := entAuth.Role{
		Name: model.Name,
		Desc: model.Desc,
	}
	result, err := a.useCase.CreateRole(ctx.Context(), user.Id, createdRole, model.PermIds...)
	if err != nil {
		return err
	}
	return ctx.JSON(roleCreateResponse{
		Id: result.Id,
	})
}

// listRoles godoc
// @ID listRoles
// @Summary list all the roles in the system
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} listRoleResponseItem
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Router /roles [get]
func (a *API) listRoles(ctx *fiber.Ctx) error {
	user, ok := api.GetUserEntity(ctx)
	if !ok {
		return errors.New("user unauthorized")
	}

	result, err := a.useCase.ListRoles(ctx.Context(), user.Id)
	if err != nil {
		return err
	}
	resp := lo.Map(result, func(r entAuth.Role, _ int) listRoleResponseItem {
		return listRoleResponseItem{
			Id:      r.Id,
			Name:    r.Name,
			Desc:    r.Desc,
			PermIds: lo.Map(r.Perms, func(p entAuth.Perm, _ int) entity.BigIntPK { return p.Id }),
		}
	})
	return ctx.JSON(resp)
}

// listPerms godoc
// @ID listPerms
// @Summary list all the perms in the system
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} listPermResponseItem
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Router /perms [get]
func (a *API) listPerms(ctx *fiber.Ctx) error {
	user, ok := api.GetUserEntity(ctx)
	if !ok {
		return errors.New("user unauthorized")
	}

	result, err := a.useCase.ListPerms(ctx.Context(), user.Id)
	if err != nil {
		return err
	}
	resp := lo.Map(result, func(p entAuth.Perm, _ int) listPermResponseItem {
		return listPermResponseItem{
			Id:   p.Id,
			Name: string(p.Name),
			Desc: p.Desc,
		}
	})
	return ctx.JSON(resp)
}

// updateRole godoc
// @ID updateRole
// @Summary update role
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} updateRoleResponse
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Param id path number true "id of a role to update"
// @Param form body updateRoleRequest true "new role data to store"
// @Security ApiKeyAuth
// @Router /roles/{id} [put]
func (a *API) updateRole(ctx *fiber.Ctx) error {
	user, ok := api.GetUserEntity(ctx)
	if !ok {
		return errors.New("user unauthorized")
	}

	roleId, err := entity.ParseBigIntPK(ctx.Params("id"))
	if err != nil {
		return err
	}

	var model updateRoleRequest
	if err := ctx.BodyParser(&model); err != nil {
		return err
	}

	_, err = a.useCase.UpdateRole(ctx.Context(), user.Id, roleId, model.Name, model.PermIds...)
	if err != nil {
		return err
	}

	return ctx.JSON(updateRoleResponse{})
}

// deleteRole godoc
// @ID deleteRole
// @Summary delete role
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} api.OK
// @Failure 400 {object} api.ErrResponse
// @Failure 500 {object} api.ErrResponse
// @Security ApiKeyAuth
// @Param id path number true "id of a role to delete"
// @Router /roles/{id} [delete]
func (a *API) deleteRole(ctx *fiber.Ctx) error {
	user, ok := api.GetUserEntity(ctx)
	if !ok {
		return errors.New("user unauthorized")
	}

	roleId, err := entity.ParseBigIntPK(ctx.Params("id"))
	if err != nil {
		return err
	}

	if err := a.useCase.DeleteRole(ctx.Context(), user.Id, roleId); err != nil {
		return err
	}

	return ctx.JSON(api.OK{})
}
