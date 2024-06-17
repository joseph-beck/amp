// Github Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Binding is used for binding data in modelling languages.
package binding

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type jsonBinding struct{}

// Get the name of the JSON binder.
func (j jsonBinding) Name() string {
	return "json"
}

// Bind a request to a given reference to any object, errors if cannot decode struct.
func (j jsonBinding) Bind(request *http.Request, obj any) error {
	buff, err := readBody(request)
	if err != nil {
		return err
	}

	return j.decodeJSON(buff, obj)
}

// Binds the a byte array to an object of type any.
func (j jsonBinding) BindBody(body []byte, obj any) error {
	return j.decodeJSON(bytes.NewReader(body), obj)
}

// Decodes a reader with an object of any, errors if it cannot validate the struct.
func (j jsonBinding) decodeJSON(reader io.Reader, obj any) error {
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(obj)
	if err != nil {
		return err
	}

	err = validate(obj)
	if err != nil {
		return err
	}

	return nil
}
