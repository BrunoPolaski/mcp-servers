package dto

type HealthResponse struct {
	APIStatus      string `json:"api_status"`
	DatabaseStatus string `json:"database_status"`
	CacheStatus    string `json:"cache_status"`
}
