package entities

// UserRole defines the type for user roles.
type UserRole string

const (
	AdminRole UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User represents a user in the system.
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}

// UserCountResponse is used for broadcasting the number of active users.
type UserCountResponse struct {
	ActiveUsers int `json:"activeUsers"`
	TotalUsers  int `json:"totalUsers"`
}
