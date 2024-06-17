// Github Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Binding is used for binding data in modelling languages.
package binding

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
)

type xmlBinding struct{}

func (x xmlBinding) Name() string {
	return "xml"
}

func (x xmlBinding) Bind(request *http.Request, obj any) error {
	buff, err := readBody(request)
	if err != nil {
		return err
	}

	return x.decodeXML(buff, obj)
}

func (x xmlBinding) BindBody(body []byte, obj any) error {
	return x.decodeXML(bytes.NewReader(body), obj)
}
func (x xmlBinding) decodeXML(reader io.Reader, obj any) error {
	decoder := xml.NewDecoder(reader)

	err := decoder.Decode(obj)
	if err != nil {
		return err
	}

	return validate(obj)
}
