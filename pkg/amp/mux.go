package amp

import (
	"errors"
	"fmt"
	"net/http"
)

const amp = `
________  _____ ______   ________   
|\   __  \|\   _ \  _   \|\   __  \  
\ \  \|\  \ \  \\\__\ \  \ \  \|\  \ 
 \ \   __  \ \  \\|__| \  \ \   ____\
  \ \  \ \  \ \  \    \ \  \ \  \___|
   \ \__\ \__\ \__\    \ \__\ \__\   
    \|__|\|__|\|__|     \|__|\|__|   
`

type Config struct {
	port uint
	host string
	crt  string
	key  string
}

func Default() Config {
	return Config{
		port: 8080,
		host: "",
		crt:  "",
		key:  "",
	}
}

type Mux struct {
	mux  *http.ServeMux
	port uint
	host string
	crt  string
	key  string
}

func New(config ...Config) Mux {
	c := Default()

	if len(config) > 0 {
		c = config[0]
	}

	return Mux{
		mux:  http.NewServeMux(),
		port: c.port,
		host: c.host,
		crt:  c.crt,
		key:  c.key,
	}
}

func (m *Mux) ListenAndServe() error {
	fmt.Print(amp)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", m.host, m.port), m.mux)
}

func (m *Mux) ListenAndServeTLS() error {
	if m.crt == "" || m.key == "" {
		return errors.New("error, no crt or key given")
	}

	fmt.Print(amp)

	return http.ListenAndServeTLS(fmt.Sprintf("%s:%d", m.host, m.port), m.crt, m.key, m.mux)
}
