package user

import (
	"context"
	"log"

	"github.com/MartinZitterkopf/gocurse_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error)
		GetAll(ctx context.Context, filters Fillters, offset, limit int) ([]domain.User, error)
		GetByID(ctx context.Context, id string) (*domain.User, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error
		Count(ctx context.Context, filters Fillters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Fillters struct {
		FirstName string
		LastName  string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error) {

	s.log.Println("create user service")
	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		s.log.Println(err)
		return nil, err
	}

	return &user, nil
}

func (s service) GetAll(ctx context.Context, filters Fillters, offset, limit int) ([]domain.User, error) {

	users, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	return users, nil
}

func (s service) GetByID(ctx context.Context, id string) (*domain.User, error) {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	return user, nil
}

func (s service) Delete(ctx context.Context, id string) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s service) Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error {

	if err := s.repo.Update(ctx, id, firstName, lastName, email, phone); err != nil {
		return err
	}

	return nil
}

func (s service) Count(ctx context.Context, filters Fillters) (int, error) {
	return s.repo.Count(ctx, filters)
}
