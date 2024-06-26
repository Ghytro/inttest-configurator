package mockservice

import (
	serviceRules "configurator/internal/businessrules/mockservice"
	"configurator/internal/entity"
	"configurator/pkg/exportstruct"
	"configurator/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pb33f/libopenapi"
	libopenapitypesv2 "github.com/pb33f/libopenapi/datamodel/high/v2"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type RestServiceIdentifier struct {
	ProjectId entity.BigIntPK
	ServiceId string
}

type restServiceUseCase struct {
	log      *zap.SugaredLogger
	projRepo IProjectRepository
}

type CreateRestServiceReq struct {
	Id   string `json:"id"`
	Port int    `json:"port"`
}

type CreateRestServiceResult entity.Mixin[entity.BigIntPK]

func (uc *restServiceUseCase) CreateRestService(
	ctx context.Context,
	projectId entity.BigIntPK,
	request CreateRestServiceReq,
) (result CreateRestServiceResult, err error) {
	err = uc.projRepo.ModifyProjectData(ctx, projectId, func(data *exportstruct.Config) error {
		if err := serviceRules.ValidateRestUniqServiceId(*data, request.Id); err != nil {
			return err
		}
		if err := serviceRules.ValidateUniqServicePort(*data, request.Port); err != nil {
			return err
		}

		data.RpcServices = append(data.RpcServices, exportstruct.RpcService{
			RpcServiceCommon: exportstruct.RpcServiceCommon{
				ID:   request.Id,
				Type: exportstruct.RpcServiceType_REST,
				Port: exportstruct.Port{
					Port: request.Port,
				},
			},
		})

		return nil
	})
	if err != nil {
		return CreateRestServiceResult{}, err
	}
	return CreateRestServiceResult{}, nil
}

type RestServiceUpdateForm struct {
	// todo: id выглядит странновато, мб переделать на name?
	// если мигрировать на норм бд, это поле может вызвать непонимание
	Id   string `json:"id"`
	Port int    `json:"port"`
}

type RestServiceUpdateResult entity.Mixin[entity.BigIntPK]

func (uc *restServiceUseCase) UpdateRestService(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	request RestServiceUpdateForm,
) (RestServiceUpdateResult, error) {
	err := uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис с идентификатором %s", serviceId.ServiceId)
		}
		if request.Id != service.ID {
			// исключим сервис из проверяемых в правиле и проверим уникальность id
			services := append(data.RpcServices[:serviceIdx], data.RpcServices[serviceIdx+1:]...)
			mockData := exportstruct.Config{
				RpcServices: services,
			}
			if err := serviceRules.ValidateRestUniqServiceId(mockData, request.Id); err != nil {
				return err
			}
			data.RpcServices[serviceIdx].ID = request.Id
		}
		if request.Port != service.GetPort() {
			data.RpcServices[serviceIdx].Port.Port = request.Port
		}
		return nil
	})
	if err != nil {
		return RestServiceUpdateResult{}, err
	}
	return RestServiceUpdateResult{}, nil
}

func (uc *restServiceUseCase) DeleteRestService(ctx context.Context, serviceId RestServiceIdentifier) error {
	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		_, serviceIdx, _ := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if serviceIdx == -1 {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		data.RpcServices = append(data.RpcServices[:serviceIdx], data.RpcServices[serviceIdx+1:]...)
		return nil
	})
}

type ListRestRoutesResult struct {
	Id     int    `json:"id"`
	Route  string `json:"route"`
	Method string `json:"method"`
}

func (uc *restServiceUseCase) ListRestHandlers(ctx context.Context, serviceId RestServiceIdentifier) ([]ListRestRoutesResult, error) {
	data, err := fetchProjectData(ctx, uc.projRepo, serviceId.ProjectId)
	if err != nil {
		return nil, err
	}
	service, ok := lo.Find(data.RpcServices, func(item exportstruct.RpcService) bool {
		return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
	})
	if !ok {
		return nil, fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
	}
	return lo.Map(service.RpcServiceUnion.HttpService.Routes, func(r exportstruct.HttpRoute, idx int) ListRestRoutesResult {
		return ListRestRoutesResult{
			Id:     idx,
			Route:  r.Route.String(),
			Method: string(r.Method),
		}
	}), nil
}

type CreateRestHandlerForm struct {
	Route  string `json:"route"`
	Method string `json:"method"`
}

func (uc *restServiceUseCase) AddRestHandler(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	form CreateRestHandlerForm,
) error {
	route, method := exportstruct.ParametrizedRestRoute(form.Route), exportstruct.HttpMethod(form.Method)
	if err := exportstruct.HttpMethod(form.Method).Validate(); err != nil {
		return err
	}
	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		routes := service.RpcServiceUnion.HttpService.Routes
		data.RpcServices[serviceIdx].RpcServiceUnion.HttpService.Routes = append(routes, exportstruct.HttpRoute{
			Route:  route,
			Method: method,
		})
		return nil
	})
}

type UpdateRestHandlerForm struct {
	Route  string `json:"route"`
	Method string `json:"method"`
}

func (uc *restServiceUseCase) UpdateRestHandler(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	id int,
	form UpdateRestHandlerForm,
) error {
	if id < 0 {
		return errors.New("некорректный идентификатор http-маршрута")
	}
	method := exportstruct.HttpMethod(form.Method)
	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		if id >= len(service.RpcServiceUnion.HttpService.Routes) {
			return errors.New("некорректный идентификатор http-маршрута")
		}
		if form.Route != service.RpcServiceUnion.HttpService.Routes[id].Route.String() {
			data.RpcServices[serviceIdx].RpcServiceUnion.HttpService.Routes[id].Route = exportstruct.ParametrizedRestRoute(form.Route)
		}
		if form.Method != string(service.RpcServiceUnion.HttpService.Routes[id].Method) {
			data.RpcServices[serviceIdx].RpcServiceUnion.HttpService.Routes[id].Method = method
		}
		return nil
	})
}

func (uc *restServiceUseCase) DeleteRestHandler(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	id int,
) error {
	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		if id >= len(service.RpcServiceUnion.HttpService.Routes) {
			return errors.New("некорректный идентификатор http-маршрута")
		}
		routes := data.RpcServices[serviceIdx].RpcServiceUnion.HttpService.Routes
		data.RpcServices[serviceIdx].RpcServiceUnion.HttpService.Routes = append(routes[:id], routes[id+1:]...)
		return nil
	})
}

type CreateRestStubBehaviorForm struct {
	Headers   map[string]string `json:"headers"`
	Query     map[string]string `json:"query"`
	UrlParams map[string]string `json:"url_params"`
	Body      string            `json:"body"`

	ResponseStatus  int               `json:"response_status"`
	ResponseHeaders map[string]string `json:"response_headers"`
	ResponseBody    string            `json:"response_body"`
}

func (uc *restServiceUseCase) AddRestStubBehavior(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	handlerId int,
	form CreateRestStubBehaviorForm,
) error {
	if form.Body != "" && !json.Valid(utils.S2B(form.Body)) {
		return errors.New("form body is not json encoded")
	}
	if form.ResponseBody != "" && !json.Valid(utils.S2B(form.ResponseBody)) {
		return errors.New("stub response is not json encoded")
	}
	if handlerId < 0 {
		return errors.New("некорректный id http-маршрута")
	}
	if err := exportstruct.HttpStatus(form.ResponseStatus).Validate(); err != nil {
		return err
	}

	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		if handlerId >= len(service.RpcServiceUnion.HttpService.Routes) {
			return errors.New("некорректный id http-маршрута")
		}
		behavior := data.RpcServices[serviceIdx].RpcServiceUnion.HttpService.Routes[handlerId].Behavior
		data.RpcServices[serviceIdx].RpcServiceUnion.HttpService.Routes[handlerId].Behavior = append(
			behavior,
			exportstruct.HttpHandlerBehaviorItem{
				Type: exportstruct.RestHandlerBehaviorType_STUB,
				HttpHandlerBehaviorUnion: exportstruct.HttpHandlerBehaviorUnion{
					HttpStubBehavior: exportstruct.HttpStubBehavior{
						Params: exportstruct.HttpStubBehaviorParams{
							Headers: form.Headers,
							Query:   form.Query,
							Url:     form.UrlParams,
							Body:    form.Body,
						},
						Response: exportstruct.HttpStubBehaviorResponse{
							Status:  exportstruct.HttpStatus(form.ResponseStatus),
							Headers: form.ResponseHeaders,
							Body:    form.ResponseBody,
						},
					},
				},
			},
		)
		return nil
	})
}

type CreateRestMockBehaviorForm struct {
	Impl []string `json:"impl"`
}

func (uc *restServiceUseCase) AddRestMockBehavior(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	handlerId int,
	form CreateRestMockBehaviorForm,
) error {
	if handlerId < 0 {
		return errors.New("некорректный идентификатор http-маршрута")
	}
	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		if handlerId >= len(service.RpcServiceUnion.HttpService.Routes) {
			return errors.New("некорректный id http-маршрута")
		}
		behavior := data.RpcServices[serviceIdx].RpcServiceUnion.HttpService.Routes[handlerId].Behavior
		data.RpcServices[serviceIdx].RpcServiceUnion.HttpService.Routes[handlerId].Behavior = append(
			behavior,
			exportstruct.HttpHandlerBehaviorItem{
				Type: exportstruct.RestHandlerBehaviorType_MOCK,
				HttpHandlerBehaviorUnion: exportstruct.HttpHandlerBehaviorUnion{
					HttpMockBehavior: exportstruct.HttpMockBehavior{
						Impl: form.Impl,
					},
				},
			},
		)
		return nil
	})
}

type (
	ListRestBehaviorResult struct {
		Stubs []ListRestBehaviorResult_Stub `json:"stubs"`
		Mocks []ListRestBehaviorResult_Mock `json:"mocks"`
	}

	ListRestBehaviorResult_Stub struct {
		Id       int `json:"id"`
		Priority int `json:"priority"`

		Headers   map[string]string `json:"headers"`
		Query     map[string]string `json:"query"`
		UrlParams map[string]string `json:"url_params"`
		Body      string            `json:"body"`

		ResponseStatus  int               `json:"response_status"`
		ResponseHeaders map[string]string `json:"response_headers"`
		ResponseBody    string            `json:"response_body"`
	}

	ListRestBehaviorResult_Mock struct {
		Id       int `json:"id"`
		Priority int `json:"priority"`

		Impl []string `json:"impl"`
	}
)

func (uc *restServiceUseCase) ListRestBehaviors(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	handlerId int,
) (result ListRestBehaviorResult, err error) {
	if handlerId < 0 {
		return result, errors.New("некорректный id http-хендлера")
	}
	data, err := fetchProjectData(ctx, uc.projRepo, serviceId.ProjectId)
	if err != nil {
		return ListRestBehaviorResult{}, err
	}
	service, ok := lo.Find(data.RpcServices, func(item exportstruct.RpcService) bool {
		return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
	})
	if !ok {
		return ListRestBehaviorResult{}, fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
	}
	if handlerId >= len(service.RpcServiceUnion.HttpService.Routes) {
		return result, errors.New("некорректный id http-хендлера")
	}
	route := service.RpcServiceUnion.HttpService.Routes[handlerId]
	for i, b := range route.Behavior {
		switch b.Type {
		case exportstruct.RestHandlerBehaviorType_STUB:
			reqBody, err := beautifyJson(b.Params.Body)
			if err != nil {
				return result, err
			}
			respBody, err := beautifyJson(b.Response.Body)
			if err != nil {
				return result, err
			}
			result.Stubs = append(result.Stubs, ListRestBehaviorResult_Stub{
				Id:              i,
				Priority:        i,
				Headers:         b.Params.Headers,
				Query:           b.Params.Query,
				UrlParams:       b.Params.Url,
				Body:            reqBody,
				ResponseStatus:  int(b.Response.Status),
				ResponseHeaders: b.Response.Headers,
				ResponseBody:    respBody,
			})
		case exportstruct.RestHandlerBehaviorType_MOCK:
			result.Mocks = append(result.Mocks, ListRestBehaviorResult_Mock{
				Id:       i,
				Priority: i,
				Impl:     b.Impl,
			})
		}
	}
	return result, nil
}

type UpdateStubBehaviorForm struct {
	Headers   map[string]string `json:"headers"`
	Query     map[string]string `json:"query"`
	UrlParams map[string]string `json:"url_params"`
	Body      string            `json:"body"`

	ResponseStatus  int               `json:"response_status"`
	ResponseHeaders map[string]string `json:"response_headers"`
	ResponseBody    string            `json:"response_body"`
}

func (uc *restServiceUseCase) UpdateRestStubBehavior(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	handlerId, behaviorId int,
	form UpdateStubBehaviorForm,
) error {
	if handlerId < 0 {
		return fmt.Errorf("некорректный идентификатор http-маршрута")
	}
	if behaviorId < 0 {
		return fmt.Errorf("некорректный идентификатор поведения хендлера")
	}

	if err := exportstruct.HttpStatus(form.ResponseStatus).Validate(); err != nil {
		return err
	}

	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		if handlerId >= len(service.RpcServiceUnion.HttpService.Routes) {
			return fmt.Errorf("некорректный идентификатор http-маршрута")
		}

		servicePtr := &data.RpcServices[serviceIdx]
		routePtr := &servicePtr.RpcServiceUnion.HttpService.Routes[handlerId]
		if behaviorId >= len(routePtr.Behavior) {
			return fmt.Errorf("некорректный идентификатор поведения хендлера")
		}

		var err error

		behaviorPtr := &routePtr.Behavior[behaviorId]
		stubBehaviorPtr := &behaviorPtr.HttpHandlerBehaviorUnion.HttpStubBehavior
		stubBehaviorPtr.Params.Headers = form.Headers
		stubBehaviorPtr.Params.Query = form.Query
		stubBehaviorPtr.Params.Url = form.UrlParams

		stubBehaviorPtr.Params.Body, err = compressJsonBeforeStore(form.Body)
		if err != nil {
			return err
		}

		stubBehaviorPtr.Response.Headers = form.ResponseHeaders
		stubBehaviorPtr.Response.Status = exportstruct.HttpStatus(form.ResponseStatus)

		stubBehaviorPtr.Response.Body, err = compressJsonBeforeStore(form.ResponseBody)
		if err != nil {
			return err
		}

		return nil
	})
}

type UpdateMockBehaviorForm struct {
	Impl []string `json:"impl"`
}

func (uc *restServiceUseCase) UpdateRestMockBehavior(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	handlerId, behaviorId int,
	form UpdateMockBehaviorForm,
) error {
	if handlerId < 0 {
		return fmt.Errorf("некорректный идентификатор http-маршрута")
	}
	if behaviorId < 0 {
		return fmt.Errorf("некорректный идентификатор поведения хендлера")
	}
	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		if handlerId >= len(service.RpcServiceUnion.HttpService.Routes) {
			return fmt.Errorf("некорректный идентификатор http-маршрута")
		}

		servicePtr := &data.RpcServices[serviceIdx]
		routePtr := &servicePtr.RpcServiceUnion.HttpService.Routes[handlerId]
		if behaviorId >= len(routePtr.Behavior) {
			return fmt.Errorf("некорректный идентификатор поведения хендлера")
		}
		behaviorPtr := &routePtr.Behavior[behaviorId]
		mockBehaviorPtr := &behaviorPtr.HttpHandlerBehaviorUnion.HttpMockBehavior
		mockBehaviorPtr.Impl = form.Impl
		return nil
	})
}

func (uc *restServiceUseCase) MoveRestBehaviorPriority(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	handlerId, behaviorId int,
	newPriority int,
) error {
	if handlerId < 0 {
		return fmt.Errorf("некорректный идентификатор http-маршрута")
	}
	if behaviorId < 0 {
		return fmt.Errorf("некорректный идентификатор поведения хендлера")
	}
	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		if handlerId >= len(service.RpcServiceUnion.HttpService.Routes) {
			return fmt.Errorf("некорректный идентификатор http-маршрута")
		}
		servicePtr := &data.RpcServices[serviceIdx]
		routePtr := &servicePtr.RpcServiceUnion.HttpService.Routes[handlerId]
		if behaviorId >= len(routePtr.Behavior) {
			return fmt.Errorf("некорректный идентификатор поведения хендлера")
		}
		routePtr.Behavior = append(
			routePtr.Behavior[:newPriority],
			append(
				routePtr.Behavior[behaviorId:behaviorId+1],
				routePtr.Behavior[newPriority+1:]...,
			)...,
		)
		return nil
	})
}

func (uc *restServiceUseCase) DeleteRestBehavior(
	ctx context.Context,
	serviceId RestServiceIdentifier,
	handlerId, behaviorId int,
) error {
	if handlerId < 0 {
		return fmt.Errorf("некорректный идентификатор http-маршрута")
	}
	if behaviorId < 0 {
		return fmt.Errorf("некорректный идентификатор поведения хендлера")
	}
	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		service, serviceIdx, ok := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if !ok {
			return fmt.Errorf("не найден REST-сервис по идентификатору %s", serviceId.ServiceId)
		}
		if handlerId >= len(service.RpcServiceUnion.HttpService.Routes) {
			return fmt.Errorf("некорректный идентификатор http-маршрута")
		}
		servicePtr := &data.RpcServices[serviceIdx]
		routePtr := &servicePtr.RpcServiceUnion.HttpService.Routes[handlerId]
		if behaviorId >= len(routePtr.Behavior) {
			return fmt.Errorf("некорректный идентификатор поведения хендлера")
		}

		routePtr.Behavior = append(routePtr.Behavior[:behaviorId], routePtr.Behavior[behaviorId+1:]...)
		return nil
	})
}

func (uc *restServiceUseCase) ImportSwagger(ctx context.Context, serviceId RestServiceIdentifier, doc libopenapi.Document) error {
	if !strings.HasPrefix(doc.GetVersion(), "2") {
		return errors.New("имеется поддержка только сваггера версии 2")
	}
	model, errs := doc.BuildV2Model()
	if len(errs) != 0 {
		return &multierror.Error{
			Errors:      errs,
			ErrorFormat: multierror.ListFormatFunc,
		}
	}
	return uc.projRepo.ModifyProjectData(ctx, serviceId.ProjectId, func(data *exportstruct.Config) error {
		_, serviceIdx, _ := lo.FindIndexOf(data.RpcServices, func(item exportstruct.RpcService) bool {
			return item.Type == exportstruct.RpcServiceType_REST && item.ID == serviceId.ServiceId
		})
		if serviceIdx == -1 {
			return errors.New("incorrect rest service id")
		}
		servicePtr := &data.RpcServices[serviceIdx]
		serviceRoutes := servicePtr.RpcServiceUnion.HttpService.Routes
		for item := model.Model.Paths.PathItems.First(); item != nil; item = item.Next() {
			swagUrl := item.Key()
			methods := getSwaggerPathHttpMethods(item.Value())
			if lo.ContainsBy(serviceRoutes, func(item exportstruct.HttpRoute) bool {
				return item.Route.ToSwaggerUrl() == swagUrl && lo.Contains(methods, string(item.Method))
			}) {
				continue
			}

			for _, method := range methods {
				serviceRoutes = append(serviceRoutes, exportstruct.HttpRoute{
					Route:  swaggerToFiberUrl(swagUrl),
					Method: exportstruct.HttpMethod(method),
				})
			}
		}
		servicePtr.RpcServiceUnion.HttpService.Routes = serviceRoutes
		return nil
	})
}

func swaggerToFiberUrl(swaggerUrl string) exportstruct.ParametrizedRestRoute {
	for braceIdx := strings.Index(swaggerUrl, "{"); braceIdx != -1; braceIdx = strings.Index(swaggerUrl, "{") {
		closeBraceIdx := strings.Index(swaggerUrl, "}")
		swaggerUrl = swaggerUrl[:braceIdx] + ":" + swaggerUrl[braceIdx+1:closeBraceIdx] + swaggerUrl[closeBraceIdx+1:]
	}
	return exportstruct.ParametrizedRestRoute(swaggerUrl)
}

func getSwaggerPathHttpMethods(item *libopenapitypesv2.PathItem) (result []string) {
	pathMethods := map[*libopenapitypesv2.Operation]string{
		item.Get:     http.MethodGet,
		item.Put:     http.MethodPut,
		item.Post:    http.MethodPost,
		item.Delete:  http.MethodDelete,
		item.Options: http.MethodOptions,
		item.Head:    http.MethodHead,
		item.Patch:   http.MethodPatch,
	}
	for m, mStr := range pathMethods {
		if m != nil {
			result = append(result, mStr)
		}
	}
	return
}
