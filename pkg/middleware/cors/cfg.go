package cors

type Config struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	ExposeHeaders    []string
	Debug            bool
}
