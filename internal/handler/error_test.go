package handler

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestMissingFieldsError(t *testing.T) {
	error := &MissingFieldsError{missing: []string{"field1", "field2"}}
	assert.Equal(t, "missing required fields: field1,field2", error.Error())
}
