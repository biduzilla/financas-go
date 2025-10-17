package service

import (
	"errors"
	"financas/internal/model"
	"financas/internal/repository"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
)

type UserService struct {
	userRepository repository.UserRepository
}

type UserServiceInterface interface {
	ActivateUser(cod int, email string, v *validator.Validator) (*model.User, error)
	Update(user *model.User) error
	GetUserByCodAndEmail(cod int, email string, v *validator.Validator) (*model.User, error)
	Insert(user *model.User, v *validator.Validator) error
	RegisterUserHandler(user *model.User, v *validator.Validator) error
	GetUserByEmail(email string, v *validator.Validator) (*model.User, error)
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

func (s *UserService) GetUserByEmail(email string, v *validator.Validator) (*model.User, error) {
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ActivateUser(cod int, email string, v *validator.Validator) (*model.User, error) {
	if model.ValidateEmail(v, email); !v.Valid() {
		return nil, e.ErrInvalidData
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

func (s *UserService) GetUserByCodAndEmail(cod int, email string, v *validator.Validator) (*model.User, error) {
	user, err := s.userRepository.GetByCodAndEmail(cod, email)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrRecordNotFound):
			v.AddError("code", "invalid validation code or email")
			return nil, e.ErrInvalidData
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserService) RegisterUserHandler(user *model.User, v *validator.Validator) error {
	user.Cod = utils.GenerateRandomCode()
	return s.Insert(user, v)
}

func (s *UserService) Insert(user *model.User, v *validator.Validator) error {
	if user.ValidateUser(v); !v.Valid() {
		return e.ErrInvalidData
	}

	err := s.userRepository.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			return e.ErrInvalidData
		default:
			return err
		}
	}
	return nil
}
