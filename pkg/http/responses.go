package http

type HealthResponse struct {
	SqlStatus string `json:"sql"`
	RedisStatus string `json:"redis"`
}
