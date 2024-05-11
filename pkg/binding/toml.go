package binding

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pelletier/go-toml/v2"
)

type tomlBinding struct{}

func (t tomlBinding) Name() string {
	return "toml"
}

func (t tomlBinding) Bind(request *http.Request, obj any) error {
	buff, err := readBody(request)
	if err != nil {
		return err
	}

	return t.decodeToml(buff, obj)
}

func (t tomlBinding) BindBody(body []byte, obj any) error {
	return t.decodeToml(bytes.NewReader(body), obj)
}

func (t tomlBinding) decodeToml(reader io.Reader, obj any) error {
	decoder := toml.NewDecoder(reader)

	err := decoder.Decode(obj)
	if err != nil {
		return err
	}

	return decoder.Decode(obj)
}
