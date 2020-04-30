package oauth

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

type stringer struct{}

func (*stringer) String() string {
	return "stringer"
}

func TestParams(t *testing.T) {
	params := Params{}

	t.Run("AddFirstTime", func(t *testing.T) {
		expected := Params{"simple": {"value"}}

		params.Add("simple", "value")

		assert.Equal(t, expected, params)
	})

	t.Run("AddMultiple", func(t *testing.T) {
		expected := Params{"simple": {"value"}, "multiple": {"values", "values2", "values3"}}

		params.Add("multiple", "values", "values2")
		params.Add("multiple", "values3")

		assert.Equal(t, expected, params)
	})

	t.Run("AddStringer", func(t *testing.T) {
		expected := Params{"simple": {"value"}, "multiple": {"values", "values2", "values3"}, "stringer": {"stringer"}}

		params.Add("stringer", &stringer{})

		assert.Equal(t, expected, params)
	})

	t.Run("Set", func(t *testing.T) {
		expected := Params{"simple": {"value"}, "multiple": {"single"}, "stringer": {"stringer"}}

		params.Set("multiple", "single")

		assert.Equal(t, expected, params)
	})

}
