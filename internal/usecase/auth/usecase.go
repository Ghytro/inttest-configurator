package auth

import (
	"context"

	rulesAuth "configurator/internal/businessrules/auth"
	"configurator/internal/entity"
	entAuth "configurator/internal/entity/auth"
	"configurator/pkg/secrets"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UseCase struct {
	userRepo UserRepository
	log      *zap.SugaredLogger
}

func New(log *zap.SugaredLogger, userRepo UserRepository) *UseCase {
	return &UseCase{
		userRepo: userRepo,
		log:      log,
	}
}

func (uc *UseCase) CreateUser(
	ctx context.Context,
	creatorId entity.BigIntPK,
	userName string,
	password string,
	roles ...entity.BigIntPK,
) (newUser entAuth.User, err error) {
	err = uc.userRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		creator, err := uc.userRepo.GetUser(ctx, creatorId, true)
		if err != nil {
			uc.log.Error(err)
			return errors.New("невозможно получить данные об администраторе системы")
		}
		if err := rulesAuth.UserHasPerms(creator, entAuth.PermissionCreateUser); err != nil {
			uc.log.Error(err)
			return err
		}
		newUser = entAuth.User{
			UserName: userName,
			Password: password,
		}
		newUser, err = uc.userRepo.CreateUser(ctx, newUser, roles...)
		if err != nil {
			uc.log.Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		return entAuth.User{}, err
	}
	return newUser, nil
}

func (uc *UseCase) UpdateUser(ctx context.Context, updaterId, id entity.BigIntPK, newUsername string, newRoleIds ...entity.BigIntPK) (result entAuth.User, err error) {
	err = uc.userRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		updater, err := uc.userRepo.GetUser(ctx, updaterId, true)
		if err != nil {
			return err
		}
		if err := rulesAuth.UserHasPerms(updater, entAuth.PermissionEditUser); err != nil {
			return err
		}

		user, err := uc.userRepo.GetUser(ctx, id, false)
		if err != nil {
			return err
		}
		user.UserName = newUsername
		result, err = uc.userRepo.UpdateUser(ctx, user)
		if err != nil {
			return err
		}
		if len(newRoleIds) != 0 {
			result, err = uc.userRepo.UpdateUserRoles(ctx, result, newRoleIds...)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return entAuth.User{}, err
	}
	return result, nil
}

func (uc *UseCase) DeleteUser(
	ctx context.Context,
	deleterId entity.BigIntPK,
	userId entity.BigIntPK,
) error {
	return uc.userRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		deleter, err := uc.userRepo.GetUser(ctx, deleterId, true)
		if err != nil {
			uc.log.Error(err)
			return errors.New("ошибка получения данных об администраторе системы")
		}
		if err := rulesAuth.UserHasPerms(deleter, entAuth.PermissionDeleteUser); err != nil {
			uc.log.Error(err)
			return err
		}
		if err := uc.userRepo.DeleteUser(ctx, userId); err != nil {
			uc.log.Error(err)
			return err
		}
		return nil
	})
}

func (uc *UseCase) ListUsers(
	ctx context.Context,
	querierId entity.BigIntPK,
) (result []entAuth.User, err error) {
	err = uc.userRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		deleter, err := uc.userRepo.GetUser(ctx, querierId, true)
		if err != nil {
			uc.log.Error(err)
			return errors.New("ошибка получения данных об администраторе системы")
		}
		if err := rulesAuth.UserHasPerms(deleter, entAuth.PermissionListUser); err != nil {
			uc.log.Error(err)
			return err
		}
		result, err = uc.userRepo.ListUsers(ctx, true)
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

func (uc *UseCase) AuthUser(ctx context.Context, username, password string) (user entAuth.User, err error) {
	return uc.userRepo.AuthUser(ctx, username, password)
}

func (uc *UseCase) GenToken(ctx context.Context, username, password string) (token string, err error) {
	if err := uc.userRepo.IsAuth(ctx, username, password); err != nil {
		uc.log.Error(err)
		return "", err
	}
	claims := jwt.MapClaims{
		ClaimsKeyUsername: username,
		ClaimsKeyPassword: password,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secrets.JwtSecret)
}

func (uc *UseCase) CreateRole(ctx context.Context, creatorId entity.BigIntPK, role entAuth.Role, permIds ...entity.BigIntPK) (result entAuth.Role, err error) {
	err = uc.userRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		creator, err := uc.userRepo.GetUser(ctx, creatorId, true)
		if err != nil {
			return err
		}
		if err := rulesAuth.UserHasPerms(creator, entAuth.PermissionEditUser); err != nil {
			return err
		}

		result, err = uc.userRepo.CreateRole(ctx, role, permIds...)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entAuth.Role{}, err
	}
	return result, nil
}

func (uc *UseCase) ListRoles(ctx context.Context, querierId entity.BigIntPK) (result []entAuth.Role, err error) {
	err = uc.userRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		querier, err := uc.userRepo.GetUser(ctx, querierId, true)
		if err != nil {
			return err
		}
		if err := rulesAuth.UserHasPerms(querier, entAuth.PermissionListUser); err != nil {
			return err
		}

		result, err = uc.userRepo.ListRoles(ctx, true)
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

func (uc *UseCase) ListPerms(ctx context.Context, querierId entity.BigIntPK) (result []entAuth.Perm, err error) {
	err = uc.userRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		querier, err := uc.userRepo.GetUser(ctx, querierId, true)
		if err != nil {
			return err
		}
		if err := rulesAuth.UserHasPerms(querier, entAuth.PermissionListUser); err != nil {
			return err
		}

		result, err = uc.userRepo.ListPerms(ctx)
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

func (uc *UseCase) UpdateRole(ctx context.Context, updaterId entity.BigIntPK, roleId entity.BigIntPK, newName string, newDesc *string, newPermIds ...entity.BigIntPK) (result entAuth.Role, err error) {
	err = uc.userRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		updater, err := uc.userRepo.GetUser(ctx, updaterId, true)
		if err != nil {
			return err
		}
		if err := rulesAuth.UserHasPerms(updater, entAuth.PermissionEditUser); err != nil {
			return err
		}

		result, err = uc.userRepo.UpdateRole(ctx, entAuth.Role{Id: roleId, Name: newName, Desc: newDesc})
		if err != nil {
			return err
		}
		if len(newPermIds) != 0 {
			result, err = uc.userRepo.UpdateRolePerms(ctx, result, newPermIds...)
			if err != nil {
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

func (uc *UseCase) DeleteRole(ctx context.Context, deleterId entity.BigIntPK, id entity.BigIntPK) error {
	return uc.userRepo.RunInTransaction(ctx, func(ctx context.Context) error {
		deleter, err := uc.userRepo.GetUser(ctx, deleterId, true)
		if err != nil {
			return err
		}
		if err := rulesAuth.UserHasPerms(deleter, entAuth.PermissionDeleteUser); err != nil {
			return err
		}

		if err := uc.userRepo.DeleteRole(ctx, id); err != nil {
			return err
		}
		return nil
	})
}

func (uc *UseCase) CheckUserPermissions(ctx context.Context, id entity.BigIntPK, neededPerms ...entAuth.EPermission) error {
	user, err := uc.userRepo.GetUser(ctx, id, true)
	if err != nil {
		return err
	}
	if err := rulesAuth.UserHasPerms(user, neededPerms...); err != nil {
		return err
	}
	return nil
}
