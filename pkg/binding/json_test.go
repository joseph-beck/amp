package binding

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONBindingName(t *testing.T) {
	binder := jsonBinding{}
	assert.Equal(t, "json", binder.Name())
}

func TestJSONBindingBind(t *testing.T) {
	data := []byte(`{"key": "value"}`)
	request, err := http.NewRequest("GET", "/test/one", bytes.NewReader(data))
	assert.NoError(t, err)

	var obj Mock
	binder := jsonBinding{}

	err = binder.Bind(request, &obj)
	assert.NoError(t, err)

	data = []byte(`{"key": ""}`)
	request, err = http.NewRequest("GET", "/test/two", bytes.NewReader(data))
	assert.NoError(t, err)

	err = binder.Bind(request, &obj)
	assert.Error(t, err)

	data = []byte(`{"value": ""}`)
	request, err = http.NewRequest("GET", "/test/three", bytes.NewReader(data))
	assert.NoError(t, err)

	err = binder.Bind(request, &obj)
	assert.Error(t, err)
}

func TestJSONBindingBindBody(t *testing.T) {
	data := []byte(`{"key": "value"}`)

	var obj Mock
	binder := jsonBinding{}

	err := binder.BindBody(data, &obj)
	assert.NoError(t, err)

	data = []byte(`{"key": ""}`)

	err = binder.BindBody(data, &obj)
	assert.Error(t, err)

	data = []byte(`{"value": ""}`)

	err = binder.BindBody(data, &obj)
	assert.Error(t, err)
}
