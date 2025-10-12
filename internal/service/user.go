package service

import (
	"errors"
	"financas/internal/model"
	"financas/internal/repository"
	"financas/utils"
	"financas/utils/validator"
)

type UserService struct {
	userRepository repository.UserRepository
	Validator      *validator.Validator
}

type UserServiceInterface interface {
	ActivateUser(cod int, email string) (*model.User, error)
	Update(user *model.User) error
	GetUserByCodAndEmail(cod int, email string) (*model.User, error)
	Insert(user *model.User) error
	RegisterUserHandler(user *model.User) error
}

func NewUserService(repo repository.UserRepository, v *validator.Validator) *UserService {
	return &UserService{
		userRepository: repo,
		Validator:      v,
	}
}

func (s *UserService) ActivateUser(cod int, email string) (*model.User, error) {
	if model.ValidateEmail(s.Validator, email); !s.Validator.Valid() {
		return nil, validator.ErrInvalidData
	}

	user, err := s.userRepository.GetByCodAndEmail(cod, email)

	if err != nil {
		return nil, err
	}

	user.Activated = true
	user.Cod = 0

	if err = s.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(user *model.User) error {
	err := s.userRepository.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByCodAndEmail(cod int, email string) (*model.User, error) {
	user, err := s.userRepository.GetByCodAndEmail(cod, email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			s.Validator.AddError("code", "invalid validation code or email")
			return nil, validator.ErrInvalidData
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserService) RegisterUserHandler(user *model.User) error {
	user.Cod = utils.GenerateRandomCode()
	return s.Insert(user)
}

func (s *UserService) Insert(user *model.User) error {
	if model.ValidateUser(s.Validator, user); s.Validator.Valid() {
		return validator.ErrInvalidData
	}

	err := s.userRepository.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			s.Validator.AddError("email", "a user with this email address already exists")
			return validator.ErrInvalidData
		default:
			return err
		}
	}
	return nil
}
