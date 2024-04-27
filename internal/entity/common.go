package entity

import (
	"strconv"
	"time"

	"github.com/go-pg/pg/v10/types"
)

type BigIntPK uint64

func ParseBigIntPK(s string) (BigIntPK, error) {
	result, err := strconv.ParseUint(s, 10, 64)
	return BigIntPK(result), err
}

type BaseTimestamps struct {
	CreatedAt time.Time      `json:"created_at" pg:"created_at"`
	UpdatedAt types.NullTime `json:"updated_at" pg:"updated_at"`
	DeletedAt types.NullTime `json:"deleted_at" pg:"deleted_at"`
}

type Mixin[TId any] struct {
	Id TId `json:"id" pg:"id,pk"`

	CreatedAt time.Time      `json:"created_at" pg:"created_at"`
	UpdatedAt types.NullTime `json:"updated_at" pg:"updated_at"`
	DeletedAt types.NullTime `json:"deleted_at" pg:"deleted_at"`
}
