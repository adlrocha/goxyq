package config

// Config structure
type Config struct {
	Port            string
	DestinationHost string
	ProxyPathPrefix string
	HeaderBypass    []string
	//TODO: Allow for multiple queue attributes in the proxy. A []string
	QueueAtrribute string
}

// GetConfig returns current config.
func GetConfig() *Config {
	return &Config{
		Port:            ":9099",
		DestinationHost: "http://localhost:8080",
		ProxyPathPrefix: "/api/v1/",
		QueueAtrribute:  "token",
		HeaderBypass:    []string{"X-API-KEY", "Content-Type", "Accept"}}
}
