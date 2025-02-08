package dto

type HealthCheck struct {
	AppName string `json:"app_name"`
	Message string `json:"message,omitempty"`
	Version string `json:"version,omitempty"`
}
