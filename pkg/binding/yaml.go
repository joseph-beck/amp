// GitHub Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Binding is used for binding data in modelling languages.
package binding

import (
	"bytes"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (y yamlBinding) Bind(request *http.Request, obj any) error {
	buff, err := readBody(request)
	if err != nil {
		return err
	}

	return y.decodeYAML(buff, obj)
}

func (y yamlBinding) BindBody(body []byte, obj any) error {
	return y.decodeYAML(bytes.NewReader(body), obj)
}

func (y yamlBinding) decodeYAML(reader io.Reader, obj any) error {
	decoder := yaml.NewDecoder(reader)

	err := decoder.Decode(obj)
	if err != nil {
		return err
	}

	return validate(obj)
}
