package api

import (
	"errors"
	"fmt"
	"strconv"

	"configurator/internal/entity"
	entAuth "configurator/internal/entity/auth"

	"github.com/samber/lo"

	"github.com/gofiber/fiber/v2"
)

type ErrResponse struct {
	Msg      string `json:"message"`
	Location string `json:"location"`
	HttpCode int    `json:"-"`
}

func (err ErrResponse) Error() string {
	return err.Msg
}

type OK struct{}

func ParseBody[T any](ctx *fiber.Ctx) (T, error) {
	var form T
	if err := ctx.BodyParser(&form); err != nil {
		return lo.Empty[T](), err
	}
	return form, nil
}

type parseable interface {
	entity.BigIntPK | string | int
}

func ParseUrlParamsId[T parseable](ctx *fiber.Ctx, key string) (T, error) {
	strVal := ctx.Params(key)
	if strVal == "" {
		return lo.Empty[T](), fmt.Errorf("отсутствует значение %s в url-параметрах", key)
	}
	switch any(lo.Empty[T]()).(type) {
	case entity.BigIntPK:
		result, err := entity.ParseBigIntPK(ctx.Params(key))
		if err != nil {
			return lo.Empty[T](), err
		}
		return any(result).(T), nil
	case int:
		result, err := strconv.Atoi(ctx.Params(key))
		if err != nil {
			return lo.Empty[T](), err
		}
		return any(result).(T), nil
	case string:
		return any(ctx.Params(key)).(T), nil
	}
	return lo.Empty[T](), errors.New("unknown id type in url params")
}

func ParseQueryParam[T parseable](ctx *fiber.Ctx, query string) (T, error) {
	switch any(lo.Empty[T]()).(type) {
	case entity.BigIntPK:
		result, err := entity.ParseBigIntPK(ctx.Query(query))
		if err != nil {
			return lo.Empty[T](), err
		}
		return any(result).(T), nil
	case int:
		result, err := strconv.Atoi(ctx.Query(query))
		if err != nil {
			return lo.Empty[T](), err
		}
		return any(result).(T), nil
	case string:
		return any(ctx.Query(query)).(T), nil
	}
	return lo.Empty[T](), errors.New("unknown id type in url params")
}

type Authenticator struct {
	uc AuthMiddlewareUseCase
}

func NewAuthenticator(uc AuthMiddlewareUseCase) Authenticator {
	return Authenticator{
		uc: uc,
	}
}

func (a Authenticator) AuthMiddleware() fiber.Handler {
	return NewAuthMiddleware(a.uc)
}

func (a Authenticator) CheckPermsMiddleware(perms ...entAuth.EPermission) fiber.Handler {
	return MakePermValidator(a.uc)(perms...)
}
