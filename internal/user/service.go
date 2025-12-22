package user

import (
	"context"
	"log"

	"github.com/JuD4Mo/go_api_web_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
		Get(ctx context.Context, id string) (*domain.User, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, firstName, lastName, email, phone *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
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
	s.log.Println("Create user service")
	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}
	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	usersResult, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}

	return usersResult, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.User, error) {
	userResult, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return userResult, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
func (s service) Update(ctx context.Context, id string, firstName, lastName, email, phone *string) error {
	return s.repo.Update(ctx, id, firstName, lastName, email, phone)
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
