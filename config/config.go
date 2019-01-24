package config

// Config structure
type Config struct {
	Port            string
	DestinationHost string
	HeaderBypass    []string
}

// GetConfig returns current config.
func GetConfig() *Config {
	return &Config{
		Port:            ":9090",
		DestinationHost: "https://tokenapi.tid.es",
		HeaderBypass:    []string{"X-API-KEY", "Content-Type", "Accept"}}
}
