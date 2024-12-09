package config

type Config struct {
	Host        string
	Port        string
	DatabaseUrl string
	JwtSecret   string
}

func ReadConfig(getenv func(string) string) (*Config, error) {
	config := &Config{}
	config.Host = getenv("Host")
	config.Port = getenv("Port")
	config.DatabaseUrl = getenv("DatabaseUrl")
	config.JwtSecret = getenv("JwtSecret")
	return config, nil
}
