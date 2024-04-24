package projects

import (
	"configurator/internal/entity"
	"configurator/pkg/exportstruct"
	"encoding/json"
)

type Project struct {
	Id      entity.BigIntPK `json:"id" pg:"id,pk"`
	Name    string          `json:"name" pg:"name"`
	Desc    *string         `json:"desc" pg:"description"`
	RawData json.RawMessage `json:"-" pg:"data,type:jsonb"`
}

func (p Project) ParsedData() (conf exportstruct.Config, err error) {
	err = json.Unmarshal([]byte(p.RawData), &conf)
	return
}

type ProjectListFilter struct {
}
