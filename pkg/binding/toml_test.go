package binding

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTOMLBindingName(t *testing.T) {
	binder := tomlBinding{}
	assert.Equal(t, "toml", binder.Name())
}

func TestTOMLBindingBind(t *testing.T) {
	data := []byte(`key = "value"`)
	request, err := http.NewRequest("GET", "/test/one", bytes.NewReader(data))
	assert.NoError(t, err)

	var obj Mock
	binder := tomlBinding{}

	err = binder.Bind(request, &obj)
	assert.NoError(t, err)

	data = []byte(`key = key`)
	request, err = http.NewRequest("GET", "/test/two", bytes.NewReader(data))
	assert.NoError(t, err)

	err = binder.Bind(request, &obj)
	assert.Error(t, err)

	data = []byte(`value = value`)
	request, err = http.NewRequest("GET", "/test/three", bytes.NewReader(data))
	assert.NoError(t, err)

	err = binder.Bind(request, &obj)
	assert.Error(t, err)
}

func TestTOMLBindingBindBody(t *testing.T) {
	data := []byte(`key = "value"`)

	var obj Mock
	binder := tomlBinding{}

	err := binder.BindBody(data, &obj)
	assert.NoError(t, err)

	data = []byte(`key = key`)

	err = binder.BindBody(data, &obj)
	assert.Error(t, err)

	data = []byte(`value = value`)

	err = binder.BindBody(data, &obj)
	assert.Error(t, err)

	val := make(map[string]string)
	err = binder.BindBody([]byte(`
		key = "value"
		other = "field"
	`), &val)
	assert.NoError(t, err)
	assert.Len(t, val, 2)
	assert.Equal(t, "value", val["key"])
	assert.Equal(t, "field", val["other"])
}
