package projects

import "configurator/internal/entity"

type (
	listProjectResponseItem struct {
		Id   entity.BigIntPK `json:"id"`
		Name string          `json:"name"`
		Desc *string         `json:"desc"`
	}

	updateProjectMetaReqest struct {
		Name string  `json:"name"`
		Desc *string `json:"desc"`
	}
)
