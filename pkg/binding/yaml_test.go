package binding

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLBindingName(t *testing.T) {
	binder := tomlBinding{}
	assert.Equal(t, "toml", binder.Name())
}

func TestYAMLBindingBind(t *testing.T) {
	data := []byte(`key: value`)
	request, err := http.NewRequest("GET", "/test/one", bytes.NewReader(data))
	assert.NoError(t, err)

	var obj Mock
	binder := yamlBinding{}

	err = binder.Bind(request, &obj)
	assert.NoError(t, err)

	data = []byte(`key`)
	request, err = http.NewRequest("GET", "/test/two", bytes.NewReader(data))
	assert.NoError(t, err)

	err = binder.Bind(request, &obj)
	assert.Error(t, err)
}

func TestYAMLBindingBindBody(t *testing.T) {
	data := []byte(`key: value`)

	var obj Mock
	binder := yamlBinding{}

	err := binder.BindBody(data, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Key, "value")

	data = []byte(`key`)

	err = binder.BindBody(data, &obj)
	assert.Error(t, err)

	val := make(map[string]string)
	err = binder.BindBody([]byte(`key: value`), &val)
	assert.NoError(t, err)
	assert.Len(t, val, 1)
	assert.Equal(t, "value", val["key"])
}
