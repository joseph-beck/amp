// Github Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Amp is a web framework made using the Go 1.22 Mux.
// Please ensure you are using Go 1.22, minimum, when using Amp.
package amp

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

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
// Otherwise uses the default configuration.
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

// Adds the default options for the Mux.
// Mostly used with cors pre-flight checks.
func (m *Mux) addOptions() {
	m.Options("/*", func(ctx *Ctx) error {
		ctx.Status(status.OK)
		return nil
	})
}

// Makes a standard net/http HandlerFunc from a handler and middleware,
// this is used when adding a given method to the Mux.
// Makes a Ctx within that contains the handlers from both the current middleware and given middleware.
// Handles any errors that occur within each handler and any aborts that occur in handlers.
// SLOGS info about the Handler also.
func (m *Mux) Make(handler Handler, middleware ...Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx(w, r)

		// constructs the func slice for the ctx, the is iterated on.
		ctx.handlers = append(ctx.handlers, m.middleware...)
		ctx.handlers = append(ctx.handlers, middleware...)
		ctx.handlers = append(ctx.handlers, handler)

		// using iteration is preferred here.
		// can still use Next for recursive Handler calling, less optimal.
		for i, h := range ctx.handlers {
			// we can skip within our handlers, if the index does not match, lets keep iterating.
			ctx.index++
			if ctx.index != i {
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

// Add middleware to the Mux.
// Uses a variadic variable here, can give one or many middlewares here.
// If there are none specified it will just return.
//
//	a.Use(func(ctx *amp.Ctx) error {
//		// define your middleware logic here
//		return ctx.Next() // can return ctx.Next()
//	})
//
//	a.Use(func(ctx *amp.Ctx) error {
//		// define your middleware logic here
//		return ctx.Next() // can return ctx.Next()
//	}, func(ctx *amp.Ctx) error {
//		// define more middleware logic here
//		err := result()
//		if err != nil {
//			// all errors handled by the Mux.
//			return err
//		}
//
//		// can also abort the Ctx, ending it at the end of this method.
//		// abort is however ignored when you use ctx.Next().
//		ctx.Abort()
//
//		return nil // can also use nil, Mux will iterate through methods.
//	})
//
// This will be applied to all routes added to the Mux after this statement is called.
func (m *Mux) Use(middleware ...Handler) {
	if len(middleware) <= 0 {
		return
	}

	m.middleware = append(m.middleware, middleware...)
}

// Generic handler, this can be used for a variety of http methods unlike specified ones, like Get.
// All given middleware will only be applied to this route.
// Will likely have to use a switch case statement within the handler to specify method.
//
//	func(ctx *amp.Ctx) error {
//		switch ctx.Method() {
//		case "GET":
//			...
//		default:
//			...
//		}
//	}
//
// Generally recommended to use a specified method.
func (m *Mux) Handler(path string, handler Handler, middleware ...Handler) {
	slog.Info("HANDLER " + path)
	m.mux.HandleFunc(path, m.Make(handler, middleware...))
}

// Create a Get route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Get requests should be used to retrieve data.
func (m *Mux) Get(path string, handler Handler, middleware ...Handler) {
	slog.Info("GET " + path)
	m.mux.HandleFunc(fmt.Sprintf("GET %s", path), m.Make(handler, middleware...))
}

// Create a Post route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Post methods should be used for posting data or changing state.
func (m *Mux) Post(path string, handler Handler, middleware ...Handler) {
	slog.Info("POST " + path)
	m.mux.HandleFunc(fmt.Sprintf("POST %s", path), m.Make(handler, middleware...))
}

// Create a Put route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Put methods should be used for posting data, changing state or updating state.
func (m *Mux) Put(path string, handler Handler, middleware ...Handler) {
	slog.Info("PUT " + path)
	m.mux.HandleFunc(fmt.Sprintf("PUT %s", path), m.Make(handler, middleware...))
}

// Create a Patch route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Patch methods should be used for changing state or updating data.
func (m *Mux) Patch(path string, handler Handler, middleware ...Handler) {
	slog.Info("PATCH " + path)
	m.mux.HandleFunc(fmt.Sprintf("PATCH %s", path), m.Make(handler, middleware...))
}

// Create a Delete route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Delete methods should be used for deleting data or a piece of state.
func (m *Mux) Delete(path string, handler Handler, middleware ...Handler) {
	slog.Info("DELETE " + path)
	m.mux.HandleFunc(fmt.Sprintf("DELETE %s", path), m.Make(handler, middleware...))
}

// Create a Head route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
func (m *Mux) Head(path string, handler Handler, middleware ...Handler) {
	slog.Info("HEAD " + path)
	m.mux.HandleFunc(fmt.Sprintf("HEAD %s", path), m.Make(handler, middleware...))
}

// Create an Options route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
func (m *Mux) Options(path string, handler Handler, middleware ...Handler) {
	slog.Info("OPTIONS " + path)
	m.mux.HandleFunc(fmt.Sprintf("OPTIONS %s", path), m.Make(handler, middleware...))
}

// Create a Connect route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
func (m *Mux) Connect(path string, handler Handler, middleware ...Handler) {
	slog.Info("CONNECT " + path)
	m.mux.HandleFunc(fmt.Sprintf("CONNECT %s", path), m.Make(handler, middleware...))
}

// Create a Trace route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
func (m *Mux) Trace(path string, handler Handler, middleware ...Handler) {
	slog.Info("TRACE " + path)
	m.mux.HandleFunc(fmt.Sprintf("TRACE %s", path), m.Make(handler, middleware...))
}

// Serve HTTP for the given writer and request.
// Serves the current instance of routes in Mux.
//
//	func TestXxx(t *testing.T) {
//		amp := New()
//
//		amp.Get("/test", func(ctx *Ctx) error {
//
//			...
//
//			return nil
//		})
//
//		request := httptest.NewRequest("GET", "/test", nil)
//		writer := httptest.NewRecorder()
//		amp.ServeHTTP(writer, request)
//	}
//
// More commonly used when testing routes.
func (m *Mux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.mux.ServeHTTP(writer, request)
}

// Serve your Mux one all routes have and middleware have been added.
// This will run indefinitely.
//
//	func main() {
//		a := amp.New()
//
//		a.Get("/path", func(ctx *amp.Ctx) {
//			return ctx.Render(status.OK, "hello")
//		})
//
//		log.Fatalln(a.ListenAndServe())
//	}
//
// If configured to add options, will add a catch all options method for pre-flight,
// which is mostly used for cors features.
// Can be disabled with New() and a custom configuration.
func (m *Mux) ListenAndServe() error {
	if m.defaultOptions {
		m.addOptions()
	}

	fmt.Print(amp + "\n")

	return http.ListenAndServe(fmt.Sprintf("%s:%d", m.host, m.port), m.mux)
}

// Serve your Mux one all routes have and middleware have been added.
// Please provide your configuration with a CRT and Key in order to run TLS,
// without these an error will be returned and the program likely exited.
// Allows for HTTPS requests, for better security.
// This will run indefinitely.
//
//	func main() {
//		a := amp.New(Config{
//			Port:           8080,
//			Host:           "",
//			CRT:            "CRT",
//			Key:            "Key",
//			DefaultOptions: true,
//		})
//
//		a.Get("/path", func(ctx *amp.Ctx) {
//			return ctx.Render(status.OK, "hello")
//		})
//
//		log.Fatalln(a.ListenAndServeTLS())
//	}
//
// If configured to add options, will add a catch all options method for pre-flight,
// which is mostly used for cors features.
// Can be disabled with New() and a custom configuration.
func (m *Mux) ListenAndServeTLS() error {
	if m.defaultOptions {
		m.addOptions()
	}

	if m.crt == "" || m.key == "" {
		return errors.New("error, no crt or key given")
	}

	fmt.Print(amp + "\n")

	return http.ListenAndServeTLS(fmt.Sprintf("%s:%d", m.host, m.port), m.crt, m.key, m.mux)
}
