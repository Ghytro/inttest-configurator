package api

import (
	"context"
	"errors"

	"configurator/internal/entity"
	entAuth "configurator/internal/entity/auth"

	"configurator/internal/usecase/auth"

	"configurator/pkg/secrets"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware(useCase AuthMiddlewareUseCase) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(secrets.JwtSecret),
		},
		TokenLookup: JwtLocation,
		ContextKey:  JwtContextKey,
		SuccessHandler: func(c *fiber.Ctx) error {
			return JwtValidator(useCase, c)
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "missing or malformed JWT" {
				// ErrResponse
				return c.Status(fiber.StatusUnauthorized).SendString("Отсутствует или неправильно сформирован Токен Авторизации")
			} else {
				// ErrResponse
				return c.Status(fiber.StatusUnauthorized).SendString("Недействительный или просроченный Токен Авторизации")
			}
		},
	})
}

type AuthMiddlewareUseCase interface {
	AuthUser(ctx context.Context, username, password string) (user entAuth.User, err error)
	CheckUserPermissions(ctx context.Context, id entity.BigIntPK, neededPerms ...entAuth.EPermission) error
}

func JwtValidator(useCase AuthMiddlewareUseCase, c *fiber.Ctx) error {
	// если удачно распарсили, надо сверить данные с акутальными
	token, ok := GetJwt(c)
	if !ok {
		return errors.New("jwt-token не сохранен")
	}
	claims := token.Claims.(jwt.MapClaims)
	if claims == nil {
		return errors.New("полезная нагрузка токена пуста")
	}
	username, ok := claims[auth.ClaimsKeyUsername]
	if !ok {
		return errors.New("отсутствует имя пользователя")
	}
	password, ok := claims[auth.ClaimsKeyPassword]
	if !ok {
		return errors.New("отсутствует пароль")
	}
	user, err := useCase.AuthUser(c.Context(), username.(string), password.(string))
	if err != nil {
		return err
	}
	c.Locals(UserEntityContextKey, user)
	return c.Next()
}

func GetJwt(ctx *fiber.Ctx) (*jwt.Token, bool) {
	v := ctx.Locals(JwtContextKey)
	if v == nil {
		return nil, false
	}
	return v.(*jwt.Token), true
}

func GetUserEntity(ctx *fiber.Ctx) (entAuth.User, bool) {
	v := ctx.Locals(UserEntityContextKey)
	if v == nil {
		return entAuth.User{}, false
	}
	return v.(entAuth.User), true
}

func MakePermValidator(uc AuthMiddlewareUseCase) func(perms ...entAuth.EPermission) fiber.Handler {
	return func(perms ...entAuth.EPermission) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			user, ok := GetUserEntity(ctx)
			if !ok {
				return errors.New("incorrect middleware order: fetch the user entity first")
			}
			if err := uc.CheckUserPermissions(ctx.Context(), user.Id, perms...); err != nil {
				return err
			}
			return ctx.Next()
		}
	}
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse{
		Msg:      err.Error(),
		Location: "somewhere",
	})
}
