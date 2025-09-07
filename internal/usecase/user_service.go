package usecase

import "github.com/ialekseychuk/my-place-identity/internal/domain"


type UserSevice struct {
	ur domain.UserRepository
}

func NewUserService(ur domain.UserRepository) *UserSevice {
	return &UserSevice{
		ur: ur,
	}
}

