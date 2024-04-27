package mockservice

import (
	"configurator/internal/usecase/mockservice"
)

type listServiceResponse map[string][]listServiceResult

type listServiceResult struct {
	Id   string `json:"id"`
	Port int    `json:"port"`
}

func listServiceRespByUseCaseResult(ucRes []mockservice.CommonServiceResult) listServiceResponse {
	result := listServiceResponse{}
	for _, item := range ucRes {
		result[item.Type] = append(result[item.Type], listServiceResult{
			Id:   item.Id,
			Port: item.Port,
		})
	}
	return result
}
