package auth

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"configurator/internal/entity"
	entAuth "configurator/internal/entity/auth"
	"configurator/internal/repository/internal/common"
	"configurator/pkg/database"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Repository struct {
	common.Mixin
}

func New(db *database.PGDB, log *zap.SugaredLogger) *Repository {
	return &Repository{
		Mixin: common.Mixin{
			DB:          db,
			Log:         log.With(zap.String("repo", "auth_repo")),
			ErrWrapDesc: "err in auth_repo",
		},
	}
}

func (r *Repository) GetUser(ctx context.Context, id entity.BigIntPK, deepFetchRoles bool) (result entAuth.User, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		q := db.Model(&result)
		if deepFetchRoles {
			q = q.Relation("Roles")
		}
		result.Id = id
		if err := q.WherePK().Select(); err != nil {
			return err
		}
		if deepFetchRoles {
			for i, role := range result.Roles {
				p, err := r.GetRolePerms(ctx, role.Id)
				if err != nil {
					return err
				}
				result.Roles[i].Perms = p
			}
		}
		return nil
	})
	if err != nil {
		return entAuth.User{}, r.WrapErr(err)
	}
	return result, nil
}

func (r *Repository) ListUsers(ctx context.Context, fetchRoles bool) (result []entAuth.User, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		q := db.Model(&result)
		if fetchRoles {
			q.Relation("Roles")
		}
		return q.Select()
	})
	if err != nil {
		return nil, r.WrapErr(err)
	}
	return result, nil
}

func (r *Repository) CreateUser(ctx context.Context, user entAuth.User, assignedRoleIds ...entity.BigIntPK) (entAuth.User, error) {
	const errDesc = "ошибка создания пользователя"

	assignedRoleIds = lo.Uniq(assignedRoleIds)

	err := r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		if err := r.createUserPreconds(ctx, user, assignedRoleIds...); err != nil {
			return errors.Wrap(err, errDesc)
		}

		user.CreatedAt = time.Now()
		user.UpdatedAt = types.NullTime{}
		user.DeletedAt = types.NullTime{}
		if _, err := db.Model(&user).Returning("*").Insert(); err != nil {
			r.Log.Error(err)
			return errors.New("попробуйте позже")
		}
		if _, err := db.Model(&user).
			Set("password=crypt(?, gen_salt('bf'))", user.Password).
			WherePK().
			Update(); err != nil {
			r.Log.Error(err)
			return errors.New("попробуйте позже")
		}

		if len(assignedRoleIds) > 0 {
			userRoles := lo.Map(assignedRoleIds, func(roleId entity.BigIntPK, _ int) *entAuth.UserRole {
				return &entAuth.UserRole{
					UserId: user.Id,
					RoleId: roleId,
				}
			})
			if _, err := db.Model(&userRoles).Insert(); err != nil {
				r.Log.Error(err)
				return errors.New("ошибка записи ролей")
			}
		}
		return nil
	})
	if err != nil {
		return entAuth.User{}, errors.Wrap(err, "ошибка создания пользователя")
	}
	return user, nil
}

func (r *Repository) createUserPreconds(ctx context.Context, user entAuth.User, assignedRoleIds ...entity.BigIntPK) error {
	return r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		var users []*entAuth.User
		if err := db.Model(&users).Where("username=?", user.UserName).Select(); err != nil {
			return err
		}
		if len(users) != 0 {
			return errors.New("имя пользователя обязано быть уникальным в системе")
		}

		if len(assignedRoleIds) > 0 {
			rolesCount, err := db.Model((*entAuth.Role)(nil)).WhereIn("id IN (?)", assignedRoleIds).Count()
			if err != nil {
				r.Log.Error(err)
				return errors.New("получение списка ролей")
			}
			if rolesCount != len(assignedRoleIds) {
				return errors.New("имеются некорректные значения назначаемых ролей")
			}
		}
		return nil
	})
}

func (r *Repository) IsAuth(ctx context.Context, username, password string) error {
	return r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		c, err := db.Model((*entAuth.User)(nil)).
			Where("username=? AND password=crypt(?, password)", username, password).
			Count()
		if err != nil {
			return err
		}
		if c == 0 {
			return errors.New("неправильный логин или пароль")
		}
		return nil
	})
}

func (r *Repository) AuthUser(ctx context.Context, username, password string) (result entAuth.User, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		err := db.Model(&result).
			Where("username=? AND password=crypt(?, password)", username, password).
			Select()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entAuth.User{}, nil
	}
	return result, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id entity.BigIntPK) error {
	return r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		if _, err := db.Model((*entAuth.UserRole)(nil)).Where("user_id = ?", id).Delete(); err != nil {
			return err
		}
		if _, err := db.Model(&entAuth.User{Id: id}).WherePK().Delete(); err != nil {
			return err
		}
		return nil
	})
}

func (r *Repository) GetRolePerms(ctx context.Context, roleId entity.BigIntPK) (result []entAuth.Perm, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		role := entAuth.Role{
			Id: roleId,
		}
		if err := db.Model(&role).Relation("Perms").WherePK().Select(); err != nil {
			return err
		}
		result = role.Perms
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) UpdateUser(ctx context.Context, model entAuth.User) (result entAuth.User, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		result = model
		if _, err := db.Model(&result).WherePK().Returning("*").Update(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entAuth.User{}, err
	}
	return result, nil
}

func (r *Repository) UpdateUserRoles(
	ctx context.Context,
	user entAuth.User,
	newRoleIDs ...entity.BigIntPK,
) (updated entAuth.User, err error) {
	newRoleIDs = lo.Uniq(newRoleIDs)
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		// посмотрим что все переданные новые идшники корректные
		roleCount, err := db.Model((*entAuth.Role)(nil)).WhereIn("id IN (?)", newRoleIDs).Count()
		if err != nil {
			return err
		}
		if roleCount != len(newRoleIDs) {
			return errors.New("переданы ошибочные идентификаторы новых ролей")
		}
		if _, err := db.Model((*entAuth.UserRole)(nil)).Where("user_id = ?", user.Id).Delete(); err != nil {
			return err
		}
		newUserRoles := lo.Map(newRoleIDs, func(roleId entity.BigIntPK, _ int) *entAuth.UserRole {
			return &entAuth.UserRole{
				UserId: user.Id,
				RoleId: roleId,
			}
		})
		if _, err := db.Model(&newUserRoles).Insert(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entAuth.User{}, err
	}
	return user, nil
}

func (r *Repository) CreateRole(ctx context.Context, role entAuth.Role, permIds ...entity.BigIntPK) (result entAuth.Role, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		result = role
		_, err := db.Model(&result).Returning("*").Insert()
		if err != nil {
			return err
		}
		if len(permIds) != 0 {
			rolePerms := lo.Map(permIds, func(permId entity.BigIntPK, _ int) *entAuth.RolePermission {
				return &entAuth.RolePermission{
					RoleId: result.Id,
					PermId: permId,
				}
			})
			if _, err := db.Model(&rolePerms).Insert(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return entAuth.Role{}, err
	}
	return result, nil
}

func (r *Repository) ListRoles(ctx context.Context, fetchPerms bool) (result []entAuth.Role, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		q := db.Model(&result)
		if fetchPerms {
			q.Relation("Perms")
		}
		if err := q.Select(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) ListPerms(ctx context.Context) (result []entAuth.Perm, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		return db.Model(&result).Select()
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) UpdateRole(ctx context.Context, role entAuth.Role) (result entAuth.Role, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		result = role
		if _, err := db.Model(&result).WherePK().Returning("*").Update(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entAuth.Role{}, err
	}
	return result, nil
}

func (r *Repository) UpdateRolePerms(ctx context.Context, role entAuth.Role, newPermIds ...entity.BigIntPK) (result entAuth.Role, err error) {
	err = r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		if _, err := db.Model((*entAuth.RolePermission)(nil)).Where("role_id = ?", role.Id).Delete(); err != nil {
			return err
		}
		rolePerms := lo.Map(newPermIds, func(permId entity.BigIntPK, _ int) *entAuth.RolePermission {
			return &entAuth.RolePermission{
				RoleId: role.Id,
				PermId: permId,
			}
		})
		if _, err := db.Model(&rolePerms).Insert(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entAuth.Role{}, err
	}
	return result, nil
}

func (r *Repository) DeleteRole(ctx context.Context, id entity.BigIntPK) error {
	return r.DB.RunInTransaction(ctx, func(db orm.DB) error {
		model := entAuth.Role{
			Id: id,
		}
		if _, err := db.Model((*entAuth.RolePermission)(nil)).Where("role_id = ?", id).Delete(); err != nil {
			return err
		}
		if _, err := db.Model((*entAuth.UserRole)(nil)).Where("role_id = ?", id).Delete(); err != nil {
			return err
		}
		if _, err := db.Model(&model).WherePK().Delete(); err != nil {
			return err
		}
		return nil
	})
}
