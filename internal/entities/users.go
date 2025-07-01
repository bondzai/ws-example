package entities

type UserCountResponse struct {
	ActiveUsers int `json:"activeUsers"`
	TotalUsers  int `json:"totalUsers"`
}
