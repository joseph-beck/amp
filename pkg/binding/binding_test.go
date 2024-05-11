package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Mock struct {
	Key string `binding:"required" json:"key" toml:"key" yaml:"key" xml:"key"`
}

func TestValidateStruct(t *testing.T) {
	// test validator

	validator := &defaultValidator{}

	err := validator.ValidateStruct(Mock{Key: "value"})
	assert.NoError(t, err)

	err = validator.ValidateStruct(Mock{})
	assert.Error(t, err)

	// test nil validator

	validator = &defaultValidator{}
	err = validator.ValidateStruct(nil)
	assert.NoError(t, err)

	// test arrays

	validator = &defaultValidator{}

	err = validator.ValidateStruct([]Mock{})
	assert.NoError(t, err)

	err = validator.ValidateStruct([]Mock{{}, {Key: "value"}})
	assert.Error(t, err)

	validationError, ok := err.(SliceValidationError)
	assert.True(t, ok)
	assert.Len(t, validationError, 1)

	assert.Error(t, validationError[0])
}

func TestEngine(t *testing.T) {
	validator := &defaultValidator{}
	engine := validator.Engine()
	assert.NotNil(t, engine)
}
