package repositories

type (
	UserRepository interface {
		IncreaseRealtimeUser() int
		DecreaseRealtimeUser() int
	}

	userRepository struct {
		totalUser int
	}
)

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) IncreaseRealtimeUser() int {
	r.totalUser++
	return r.totalUser
}

func (r *userRepository) DecreaseRealtimeUser() int {
	r.totalUser--

	if r.totalUser < 0 {
		r.totalUser = 0
	}

	return r.totalUser
}
