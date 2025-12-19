// GitHub Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Amp is a web framework made using the Go 1.22 Mux.
// Please ensure you are using Go 1.22, minimum, when using Amp.
package amp

// Group struct, used to group routes with a common prefix and middleware.
type group struct {
	// Prefix used for all routes in this group.
	// For example "/api/v1"
	prefix string

	// All handlers used in the group.
	// This internal struct store the method, path, handler and any middleware for each handler.
	handlers []handlerConfig

	// All middlewares that are applied to this group.
	// These are applied before the route specific middleware.
	middleware []Handler
}

// Internal struct to store handler configuration within a group.
type handlerConfig struct {
	// HTTP method of the handler
	// e.g. GET, POST, PUT, DELETE, etc.
	method string

	// Path of the handler.
	// This is without the group prefix.
	path string

	// The handler function itself.
	handler Handler

	// Any middleware specific to this route.
	// These are applied after the group middleware.
	middleware []Handler
}

// Create a new group with a given prefix and optional middleware.
// Uses a variadic variable here, can give one or many middlewares here.
// All routes in this group will use these middlewares.
//
//	g := amp.Group("/group")
//	g.Get("/hello", func(ctx *amp.Ctx) error {
//			return ctx.Render(200, "hello world!")
//	})
//
//	a := amp.New()
//
//	a.Group(g)
//
// Can now use the route /group/hello
func Group(prefix string, middleware ...Handler) group {
	if len(middleware) < 1 {
		middleware = make([]Handler, 0)
	}

	return group{
		prefix:     prefix,
		handlers:   make([]handlerConfig, 0),
		middleware: middleware,
	}
}

// Internal handler function, shortcut for generating a handlerConfig.
// Appends handlerConfig to the group.
func (g *group) handler(method string, path string, handler Handler, middleware ...Handler) {
	g.handlers = append(g.handlers, handlerConfig{
		method,
		path,
		handler,
		middleware,
	})
}

// Adds middleware to the group.
// Uses a variadic variable here, can give one or many middlewares here.
// All routes in this group will use these middlewares.
func (g *group) Use(middleware ...Handler) {
	g.middleware = append(g.middleware, middleware...)
}

// Get the prefix of the group.
func (g group) Prefix() string {
	return g.prefix
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
func (g *group) Handler(path string, handler Handler, middleware ...Handler) {
	g.handler("HANDLER", path, handler, middleware...)
}

// Create a Get route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Get requests should be used to retrieve data.
func (g *group) Get(path string, handler Handler, middleware ...Handler) {
	g.handler("GET", path, handler, middleware...)
}

// Create a Post route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Post methods should be used for posting data or changing state.
func (g *group) Post(path string, handler Handler, middleware ...Handler) {
	g.handler("POST", path, handler, middleware...)
}

// Create a Put route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Put methods should be used for posting data, changing state or updating state.
func (g *group) Put(path string, handler Handler, middleware ...Handler) {
	g.handler("PUT", path, handler, middleware...)
}

// Create a Patch route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Patch methods should be used for changing state or updating data.
func (g *group) Patch(path string, handler Handler, middleware ...Handler) {
	g.handler("PATCH", path, handler, middleware...)
}

// Create a Delete route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
// Delete methods should be used for deleting data or a piece of state.
func (g *group) Delete(path string, handler Handler, middleware ...Handler) {
	g.handler("DELETE", path, handler, middleware...)
}

// Create a Head route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
func (g *group) Head(path string, handler Handler, middleware ...Handler) {
	g.handler("HEAD", path, handler, middleware...)
}

// Create an Options route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
func (g *group) Options(path string, handler Handler, middleware ...Handler) {
	g.handler("OPTIONS", path, handler, middleware...)
}

// Create a Connect route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
func (g *group) Connect(path string, handler Handler, middleware ...Handler) {
	g.handler("CONNECT", path, handler, middleware...)
}

// Create a Trace route with a given path, handler and optional middleware.
// All given middleware will only be applied to this route.
func (g *group) Trace(path string, handler Handler, middleware ...Handler) {
	g.handler("TRACE", path, handler, middleware...)
}
