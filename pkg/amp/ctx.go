package amp

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
)

type Ctx struct {
	writer  http.ResponseWriter
	request *http.Request
	status  int

	values   map[string]any
	valuesMu sync.Mutex
}

// Set a value in the Ctx values map.
func (ctx *Ctx) Set(key string, val any) {
	ctx.valuesMu.Lock()
	defer ctx.valuesMu.Unlock()

	ctx.values[key] = val
}

// Get a value from the Ctx values map.
func (ctx *Ctx) Get(key string) (any, error) {
	ctx.valuesMu.Lock()
	defer ctx.valuesMu.Unlock()

	val, ok := ctx.values[key]
	if !ok {
		return nil, errors.New("error, could not find value in map")
	}

	return val, nil
}

// Get the path of the current Ctx.
func (ctx *Ctx) Path() string {
	return ctx.request.URL.Path
}

// Get the method of the current Ctx.
func (ctx *Ctx) Method() string {
	return ctx.request.Method
}

// Get a param from the Ctx, this will error if the param cannot be found.
func (ctx *Ctx) Param(key string) (string, error) {
	val := ctx.request.PathValue(key)
	if val == "" {
		return "", errors.New("error, param not found")
	}

	return val, nil
}

// Get a param of type int from the Ctx, this will error if the param cannot be found or cannot be converted.
func (ctx *Ctx) ParamInt(key string) (int, error) {
	str, err := ctx.Param(key)
	if err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Get a param of type float from the Ctx, this will error if the param cannot be found or cannot be converted.
func (ctx *Ctx) ParamFloat(key string) (float64, error) {
	str, err := ctx.Param(key)
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Get a param of type boolean from the Ctx, this will error if the param cannot be found or cannot be converted.
func (ctx *Ctx) ParamBool(key string) (bool, error) {
	str, err := ctx.Param(key)
	if err != nil {
		return false, err
	}

	val, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}

	return val, nil
}

// Get a query from the Ctx, this will return the default if it has any, otherwise it will error.
func (ctx *Ctx) Query(key string, def ...string) (string, error) {
	qry := ctx.request.URL.Query()
	if !qry.Has(key) {
		if len(def) > 0 {
			return def[0], nil
		}

		return "", errors.New("error, query not found and no default was given")
	}

	return qry.Get(key), nil
}

// Get a query of type int from the Ctx, this will return the default if it has any, otherwise it will error.
func (ctx *Ctx) QueryInt(key string, def ...int) (int, error) {
	qry := ctx.request.URL.Query()
	if !qry.Has(key) {
		if len(def) > 0 {
			return def[0], nil
		}

		return 0, errors.New("error, query not found and no default was given")
	}

	val, err := strconv.Atoi(qry.Get(key))
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Get a query of type float from the Ctx, this will return the default if it has any, otherwise it will error.
func (ctx *Ctx) QueryFloat(key string, def ...float64) (float64, error) {
	qry := ctx.request.URL.Query()
	if !qry.Has(key) {
		if len(def) > 0 {
			return def[0], nil
		}

		return 0, errors.New("error, query not found and no default was given")
	}

	val, err := strconv.ParseFloat(qry.Get(key), 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Get a query of type bool from the Ctx, this will return the default if it has any, otherwise it will error.
func (ctx *Ctx) QueryBool(key string, def ...bool) (bool, error) {
	qry := ctx.request.URL.Query()
	if !qry.Has(key) {
		if len(def) > 0 {
			return def[0], nil
		}

		return false, errors.New("error, query not found and no default was given")
	}

	val, err := strconv.ParseBool(qry.Get(key))
	if err != nil {
		return false, err
	}

	return val, nil
}

// Set the status of the current Ctx.
func (ctx *Ctx) Status(status int) {
	ctx.status = status
	ctx.writer.WriteHeader(status)
}

// Write a string to the Ctx writer.
func (ctx *Ctx) Write(body string) (int, error) {
	b, err := ctx.writer.Write([]byte(body))
	if err != nil {
		return 0, err
	}

	return b, nil
}

// Write a byte array to the Ctx writer
func (ctx *Ctx) WriteBytes(body []byte) (int, error) {
	b, err := ctx.writer.Write(body)
	if err != nil {
		return 0, err
	}

	return b, nil
}

// Render a string body with a status.
func (ctx *Ctx) Render(status int, body string) error {
	ctx.Status(status)

	_, err := ctx.Write(body)
	if err != nil {
		return err
	}

	return nil
}

// Render a byte array with a status.
func (ctx *Ctx) RenderBytes(status int, body []byte) error {
	ctx.Status(status)

	_, err := ctx.WriteBytes(body)
	if err != nil {
		return err
	}

	return nil
}
