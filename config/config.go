package config

// Config structure
type Config struct {
	Port            string
	DestinationHost string
	HeaderBypass    []string
	//TODO: Allow for multiple queue attributes in the proxy. A []string
	QueueAtrribute string
}

// GetConfig returns current config.
func GetConfig() *Config {
	return &Config{
		Port:            ":9090",
		DestinationHost: "https://tokenapi.tid.es",
		QueueAtrribute:  "token",
		HeaderBypass:    []string{"X-API-KEY", "Content-Type", "Accept"}}
}
