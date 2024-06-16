package amp

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/joseph-beck/amp/pkg/status"
)

// Pretty logo <3
const amp = `
________  _____ ______   ________   
|\   __  \|\   _ \  _   \|\   __  \  
\ \  \|\  \ \  \\\__\ \  \ \  \|\  \ 
 \ \   __  \ \  \\|__| \  \ \   ____\
  \ \  \ \  \ \  \    \ \  \ \  \___|
   \ \__\ \__\ \__\    \ \__\ \__\   
    \|__|\|__|\|__|     \|__|\|__|   
`

// Mux configs, allow you to customise the mux.
// Used as an arg in the New(args ...config) function.
type Config struct {
	// Selects the port that AMP Mux will run on.
	// The program will exit if the port is already in use.
	Port uint

	// Specify the host name the Mux that is on.
	// This field is not required when using ListenAndServe.
	Host string

	// CRT is required when using TLS / HTTPS.
	// This field is not required when using ListenAndServe.
	CRT string

	// Key is required when using TLS / HTTPS.
	// This field is not required when using ListenAndServe.
	Key string

	// Adds an OPTIONS method that works for all routes, /*.
	// This is used when doing pre-flight checks etc.
	// Please have this set to true if you want CORS policies to work.
	DefaultOptions bool
}

// Gives a default config,
// Port: 8080,
// Host : "",
// CRT: "",
// Key: "",
// DefaultOptions: true,
func Default() Config {
	return Config{
		Port:           8080,
		Host:           "",
		CRT:            "",
		Key:            "",
		DefaultOptions: true,
	}
}

// Mux, uses config to manipulate the mux.
// Use Get, Post, etc. to add new methods.
// To run use log.Fatalln(x.ListenAndServe()).
// Can use HTTPS with ListenAndServeTLS().
type Mux struct {
	// Go net/http serverMux, used as the server.
	mux *http.ServeMux

	// Port used by the Mux.
	// Program will exit if the port is already bound.
	port uint

	// Specify the host name the Mux that is on.
	// This field is not required.
	host string

	// CRT is required when using TLS / HTTPS.
	// This field is not required when using ListenAndServe.
	crt string

	// Key is required when using TLS / HTTPS.
	// This field is not required when using ListenAndServe.
	key string

	// Adds an OPTIONS method that works for all routes, /*.
	// This is used when doing pre-flight checks etc.
	// Please have this set to true if you want CORS policies to work.
	defaultOptions bool

	// Slice of Handlers used as middleware for all Handlers.
	// Will only apply to Handlers used after the x.Use(...) statement.
	middleware []Handler
}

// Returns a new Mux.
// Uses the given config, if len of args is greater than 0.
// Otherwises uses the default configuration.
func New(args ...Config) Mux {
	c := Default()

	if len(args) > 0 {
		c = args[0]
	}

	return Mux{
		mux:            http.NewServeMux(),
		port:           c.Port,
		host:           c.Host,
		crt:            c.CRT,
		key:            c.Key,
		defaultOptions: c.DefaultOptions,
		middleware:     make([]Handler, 0),
	}
}

func newCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		writer:  w,
		request: r,
		status:  status.Processing,
		aborted: false,

		values:   make(map[string]any),
		valuesMu: sync.Mutex{},

		handlers: []Handler{},
		index:    -1,
	}
}

func (m *Mux) addOptions() {
	m.Options("/*", func(ctx *Ctx) error {
		ctx.Status(status.OK)
		return nil
	})
}

func (m *Mux) Make(handler Handler, middleware ...Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx(w, r)

		// constructs the func slice for the ctx, the is iterated on.
		ctx.handlers = append(ctx.handlers, m.middleware...)
		ctx.handlers = append(ctx.handlers, middleware...)
		ctx.handlers = append(ctx.handlers, handler)

		// using iteration is preffered here.
		// can still use Next for recursive Handler calling, less optimal.
		for i, h := range ctx.handlers {
			// we can skip within our handlers, if the index does not match, lets keep iterating.
			if ctx.index+1 != i {
				continue
			}

			// checks if the handler has an error, slogs it if it does.
			err := h(ctx)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			// checks to see if the ctx was aborted.
			// will exit the rest of the handlers if aborted.
			// aborted ctx can be continued with the ctx.Next()
			if ctx.aborted {
				slog.Info(fmt.Sprintf("%s ABORTED %s %d", ctx.Method(), ctx.Path(), ctx.status))
				return
			}
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

func (m *Mux) Handler(path string, handler Handler, middleware ...Handler) {
	slog.Info("HANDLER " + path)
	m.mux.HandleFunc(path, m.Make(handler, middleware...))
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
	if m.defaultOptions {
		m.addOptions()
	}

	fmt.Print(amp)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", m.host, m.port), m.mux)
}

func (m *Mux) ListenAndServeTLS() error {
	if m.defaultOptions {
		m.addOptions()
	}

	if m.crt == "" || m.key == "" {
		return errors.New("error, no crt or key given")
	}

	fmt.Print(amp)

	return http.ListenAndServeTLS(fmt.Sprintf("%s:%d", m.host, m.port), m.crt, m.key, m.mux)
}
