package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type User struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Email     string
	Password  password
	Phone     string
	Activated bool
	Cod       int
	Version   int
	Deleted   bool
}

type UserDTO struct {
	ID    int64  `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserSaveDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type password struct {
	Plaintext *string
	Hash      []byte
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}
}

func (u *UserDTO) ToModel() *User {
	return &User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}
}

func (u *UserSaveDTO) ToModel() (*User, error) {
	user := &User{
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}

	err := user.Password.Set(u.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
