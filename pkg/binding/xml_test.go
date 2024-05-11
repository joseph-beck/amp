package binding

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLBindingName(t *testing.T) {
	binder := xmlBinding{}
	assert.Equal(t, "xml", binder.Name())
}

func TestXMLBindingBind(t *testing.T) {
	data := []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<root>
   			<key>value</key>
		</root>
	`)
	request, err := http.NewRequest("GET", "/test/one", bytes.NewReader(data))
	assert.NoError(t, err)

	var obj Mock
	binder := xmlBinding{}

	err = binder.Bind(request, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Key, "value")

	data = []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<root>
		</>
		</root>
	`)
	request, err = http.NewRequest("GET", "/test/two", bytes.NewReader(data))
	assert.NoError(t, err)

	err = binder.Bind(request, &obj)
	assert.Error(t, err)
}

func TestXMLBindingBindBody(t *testing.T) {
	data := []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<root>
   			<key>value</key>
		</root>
	`)

	var obj Mock
	binder := xmlBinding{}

	err := binder.BindBody(data, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "value", obj.Key)

	data = []byte(`key`)

	err = binder.BindBody(data, &obj)
	assert.Error(t, err)

	var val struct {
		Key string `xml:"key"`
		Field string `xml:"field"`
	}
	err = binder.BindBody([]byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<root>
			<key>value</key>
			<field>other</field>
		</root>
	`), &val)
	assert.NoError(t, err)
	assert.Equal(t, "value", val.Key)
	assert.Equal(t, "other", val.Field)
}
