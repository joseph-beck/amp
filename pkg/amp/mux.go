package amp

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
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
	mux        *http.ServeMux
	port       uint
	host       string
	crt        string
	key        string
	middleware []Handler
}

func New(config ...Config) Mux {
	c := Default()

	if len(config) > 0 {
		c = config[0]
	}

	return Mux{
		mux:        http.NewServeMux(),
		port:       c.port,
		host:       c.host,
		crt:        c.crt,
		key:        c.key,
		middleware: make([]Handler, 0),
	}
}

func newCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		writer:  w,
		request: r,

		values:   make(map[string]any),
		valuesMu: sync.Mutex{},
	}
}

func (m *Mux) Make(handler Handler, middleware ...Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx(w, r)

		if len(m.middleware) > 0 {
			for _, m := range m.middleware {
				err := m(ctx)
				if err != nil {
					slog.Error(err.Error())
					return
				}
			}
		}

		if len(middleware) > 0 {
			for _, m := range middleware {
				err := m(ctx)
				if err != nil {
					slog.Error(err.Error())
					return
				}
			}
		}

		err := handler(ctx)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		slog.Info(fmt.Sprintf("%s %s %d", ctx.Method(), ctx.Path(), ctx.status))
	}
}

func (m *Mux) Use(middleware ...Handler) {
	if len(middleware) <= 0 {
		return
	}

	m.middleware = append(m.middleware, middleware...)
}

func (m *Mux) Get(path string, handler Handler, middleware ...Handler) {
	slog.Info("GET " + path)
	m.mux.HandleFunc(fmt.Sprintf("GET %s", path), m.Make(handler, middleware...))
}

func (m *Mux) Post(path string, handler Handler, middleware ...Handler) {
	slog.Info("POST " + path)
	m.mux.HandleFunc(fmt.Sprintf("POST %s", path), m.Make(handler, middleware...))
}

func (m *Mux) Put(path string, handler Handler, middleware ...Handler) {
	slog.Info("PUT " + path)
	m.mux.HandleFunc(fmt.Sprintf("PUT %s", path), m.Make(handler, middleware...))
}

func (m *Mux) Patch(path string, handler Handler, middleware ...Handler) {
	slog.Info("PATCH " + path)
	m.mux.HandleFunc(fmt.Sprintf("PATCH %s", path), m.Make(handler, middleware...))
}

func (m *Mux) Delete(path string, handler Handler, middleware ...Handler) {
	slog.Info("DELETE " + path)
	m.mux.HandleFunc(fmt.Sprintf("DELETE %s", path), m.Make(handler, middleware...))
}

func (m *Mux) Head(path string, handler Handler, middleware ...Handler) {
	slog.Info("HEAD " + path)
	m.mux.HandleFunc(fmt.Sprintf("HEAD %s", path), m.Make(handler, middleware...))
}

func (m *Mux) Options(path string, handler Handler, middleware ...Handler) {
	slog.Info("OPTIONS " + path)
	m.mux.HandleFunc(fmt.Sprintf("OPTIONS %s", path), m.Make(handler, middleware...))
}

func (m *Mux) Connect(path string, handler Handler, middleware ...Handler) {
	slog.Info("CONNECT " + path)
	m.mux.HandleFunc(fmt.Sprintf("CONNECT %s", path), m.Make(handler, middleware...))
}

func (m *Mux) Trace(path string, handler Handler, middleware ...Handler) {
	slog.Info("TRACE " + path)
	m.mux.HandleFunc(fmt.Sprintf("TRACE %s", path), m.Make(handler, middleware...))
}

func (m *Mux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.mux.ServeHTTP(writer, request)
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
