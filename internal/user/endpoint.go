package user

import (
	"context"
	"errors"

	"github.com/MartinZitterkopf/gocurse_library_response/response"
	"github.com/MartinZitterkopf/gocurse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (response interface{}, err error)

	Endpoints struct {
		Create  Controller
		GetByID Controller
		GetAll  Controller
		Update  Controller
		Delete  Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetByIDReq struct {
		ID string
	}

	GetAllReq struct {
		FirstName string
		LastName  string
		Limit     int
		Page      int
	}

	UpdateReq struct {
		ID        string
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	DeleteReq struct {
		ID string
	}

	Config struct {
		PageLimDefault string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create:  makeCreateEnpoint(s),
		GetAll:  makeGetAllEnpoint(s, config),
		GetByID: makeGetByIDEnpoint(s),
		Update:  makeUpdateEnpoint(s),
		Delete:  makeDeleteEnpoint(s),
	}
}

func makeCreateEnpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		if req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", user, nil), nil
	}
}

func makeGetAllEnpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)

		filters := Fillters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.PageLimDefault)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", users, meta), nil
	}
}

func makeGetByIDEnpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDReq)

		user, err := s.GetByID(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", user, nil), nil
	}
}

func makeUpdateEnpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", nil, nil), nil
	}
}

func makeDeleteEnpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteReq)

		err := s.Delete(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", nil, nil), nil
	}
}
