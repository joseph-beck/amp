package cors

// Configure the CORS middleware.
// Used within the args of New().
// Can use Default() to build the default CORS config.
type Config struct {
}

func Default() Config {
	return Config{}
}
