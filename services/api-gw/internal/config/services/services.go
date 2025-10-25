package services

// ServiceConfig represents the configuration for a single service
type ServiceConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Config represents the configuration for all services
type Config struct {
	Services []ServiceConfig `json:"services"`
}
