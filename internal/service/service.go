package service

type Service struct {
	userService *UserService
}

func NewService(userService *UserService) *Service {
	return &Service{
		userService: userService,
	}
}
