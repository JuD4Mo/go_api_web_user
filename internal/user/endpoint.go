package user

import (
	"context"
	"errors"

	"github.com/JuD4Mo/go_api_web_meta/meta"
	"github.com/JuD4Mo/go_lib_response/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (response interface{}, err error)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		LastName  string `json:"first_name"`
		FirstName string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetReq struct {
		ID string
	}

	GetAllReq struct {
		FirstName string
		LastName  string
		Limit     int
		Page      int
	}

	DeleteReq struct {
		ID string
	}

	UpdateReq struct {
		ID        string
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	Config struct {
		LimitPage string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateReq)

		if req.FirstName == "" || req.LastName == "" {
			return nil, response.BadRequest("First name and last name must not be empty")
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", user, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)

		user, err := s.Get(ctx, req.ID)
		if err != nil {
			return nil, response.NotFound(err.Error())
		}

		return response.OK("success", user, nil), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		//Count
		quant, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, quant, config.LimitPage)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			if errors.As(err, &ErrorNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", users, meta), nil

	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		updateReq := request.(UpdateReq)

		//validaciones
		if updateReq.FirstName != nil && *updateReq.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if updateReq.LastName != nil && *updateReq.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		id := updateReq.ID
		err := s.Update(ctx, id, updateReq.FirstName, updateReq.LastName, updateReq.Email, updateReq.Phone)
		if err != nil {
			if errors.As(err, &ErrorNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("updated", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(DeleteReq)

		err := s.Delete(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrorNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("deleted", nil, nil), nil
	}
}
