package auth

import (
	"configurator/internal/entity"
	entAuth "configurator/internal/entity/auth"
	"context"
)

type IUseCase interface {
	CreateUser(ctx context.Context, creatorId entity.BigIntPK, userName string, password string, roles ...entity.BigIntPK) (newUser entAuth.User, err error)
	UpdateUser(ctx context.Context, updaterId entity.BigIntPK, id entity.BigIntPK, newUsername string, newRoleIds ...entity.BigIntPK) (result entAuth.User, err error)
	DeleteUser(ctx context.Context, deleterId entity.BigIntPK, userId entity.BigIntPK) error
	ListUsers(ctx context.Context, querierId entity.BigIntPK) (result []entAuth.User, err error)
	AuthUser(ctx context.Context, username, password string) (user entAuth.User, err error)
	GenToken(ctx context.Context, username, password string) (token string, err error)

	CreateRole(ctx context.Context, creatorId entity.BigIntPK, role entAuth.Role, permIds ...entity.BigIntPK) (result entAuth.Role, err error)
	ListRoles(ctx context.Context, querierId entity.BigIntPK) ([]entAuth.Role, error)
	ListPerms(ctx context.Context, querierId entity.BigIntPK) (result []entAuth.Perm, err error)
	UpdateRole(ctx context.Context, updaterId entity.BigIntPK, roleId entity.BigIntPK, newName string, newDesc *string, newPermIds ...entity.BigIntPK) (rule entAuth.Role, err error)
	DeleteRole(ctx context.Context, deleterId entity.BigIntPK, id entity.BigIntPK) error
}
