package usecases

import "api-gateway/internal/repositories"

type (
	UserUseCase interface {
		IncreaseRealtimeUser() int
		DecreaseRealtimeUser() int
	}

	userUseCase struct {
		repo repositories.UserRepository
	}
)

func NewUserUseCase(repo repositories.UserRepository) UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (u *userUseCase) IncreaseRealtimeUser() int {
	return u.repo.IncreaseRealtimeUser()
}

func (u *userUseCase) DecreaseRealtimeUser() int {
	return u.repo.DecreaseRealtimeUser()
}
