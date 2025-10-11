package service

import (
	"errors"
	"financas/internal/model"
	"financas/internal/repository"
	e "financas/utils/errors"
	"financas/utils/validator"
)

type UserService struct {
	userRepository *repository.UserRepository
	Validator      *validator.Validator
}

func NewUserService(repo *repository.UserRepository, v *validator.Validator, errResp *e.ErrorResponse) *UserService {
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
