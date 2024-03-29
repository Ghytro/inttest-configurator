package auth

import (
	"time"

	"configurator/internal/entity"
)

type (
	createUserRequest struct {
		UserName string   `json:"username"`
		Password string   `json:"password"`
		Roles    []string `json:"roles"`
	}
	createUserResponse struct {
		Id        entity.BigIntPK `json:"id"`
		CreatedAt time.Time       `json:"created_at"`
	}

	listUsersResponseItem struct {
		Id        entity.BigIntPK `json:"id"`
		UserName  string          `json:"username"`
		CreatedAt time.Time       `json:"created_at"`
		Roles     []string        `json:"roles"`
	}

	authRequest struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	authResponse struct {
		Token string `json:"token"`
	}
)
