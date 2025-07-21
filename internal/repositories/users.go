package repositories

import (
	"api-gateway/internal/entities"
	"context"
	"fmt"
)

// Sample users for mock data.
var (
	UserAlice   = &entities.User{ID: "user-1", Username: "Alice", Role: entities.AdminRole}
	UserBob     = &entities.User{ID: "user-2", Username: "Bob", Role: entities.RoleUser}
	UserCharlie = &entities.User{ID: "user-3", Username: "Charlie", Role: entities.RoleUser}
)

// UserRepository defines the interface for user data storage.
type UserRepository interface {
	// FindByID retrieves a user by their unique ID.
	FindByID(ctx context.Context, id string) (*entities.User, error)
	// FindByUsername retrieves a user by their username.
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
}

// MockUserRepository is an in-memory implementation of UserRepository for testing.
type MockUserRepository struct {
	users map[string]*entities.User
}

// NewMockUserRepository creates a new mock user repository with sample data.
func NewMockUserRepository() UserRepository {
	// Pre-populate with a few sample users.
	users := map[string]*entities.User{
		UserAlice.ID:   UserAlice,
		UserBob.ID:     UserBob,
		UserCharlie.ID: UserCharlie,
	}
	return &MockUserRepository{users: users}
}

// FindByID looks up a user by ID in the mock repository.
func (r *MockUserRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user with ID '%s' not found", id)
}

// FindByUsername looks up a user by username in the mock repository.
func (r *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with username '%s' not found", username)
}
