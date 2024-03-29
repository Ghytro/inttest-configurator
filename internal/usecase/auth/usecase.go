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
	roles ...string,
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
