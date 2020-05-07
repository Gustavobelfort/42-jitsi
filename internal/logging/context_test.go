package logging

import (
	"context"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
)

func TestContextGetFields(t *testing.T) {
	t.Run("NilContext", func(t *testing.T) {
		assert.Equal(t, logrus.Fields{}, ContextGetFields(nil))
	})

	t.Run("EmptyContext", func(t *testing.T) {
		assert.Equal(t, logrus.Fields{}, ContextGetFields(context.Background()))
	})

	t.Run("WithFields", func(t *testing.T) {
		expected := logrus.Fields{"expected": "fields", "id": 1}

		ctx := context.WithValue(context.Background(), logFieldsKey, expected)

		assert.Equal(t, expected, ContextGetFields(ctx))
	})
}

func TestContextWithFields(t *testing.T) {
	ctx := context.Background()

	t.Run("NoFields", func(t *testing.T) {
		expected := logrus.Fields{"expected": "fields", "id": 1}

		ctx = ContextWithFields(ctx, expected)

		assert.Equal(t, expected, ctx.Value(logFieldsKey))
	})

	t.Run("MergingFields", func(t *testing.T) {
		expected := logrus.Fields{"expected": "fields", "id": 2, "newfield": "here"}

		ctx = ContextWithFields(ctx, logrus.Fields{"id": 2, "newfield": "here"})

		assert.Equal(t, expected, ctx.Value(logFieldsKey))
	})

	t.Run("SimpleWithField", func(t *testing.T) {
		expected := logrus.Fields{"expected": "fields", "id": 3, "newfield": "here"}

		ctx = ContextWithField(ctx, "id", 3)

		assert.Equal(t, expected, ctx.Value(logFieldsKey))
	})
}

func TestContextLog(t *testing.T) {
	expected := logrus.Fields{"expected": "fields", "id": 3, "newfield": "here"}

	ctx := context.WithValue(context.Background(), logFieldsKey, expected)
	entry := ContextLog(ctx, logrus.StandardLogger())

	assert.Equal(t, expected, entry.Data)
}

func TestContextGetSentryCategory(t *testing.T) {
	t.Run("NilContext", func(t *testing.T) {
		assert.Equal(t, "", ContextGetSentryCategory(nil))
	})

	t.Run("EmptyContext", func(t *testing.T) {
		assert.Equal(t, "", ContextGetSentryCategory(context.Background()))
	})

	t.Run("WithCategory", func(t *testing.T) {
		expected := "category"

		ctx := context.WithValue(context.Background(), sentryCategoryKey, expected)

		assert.Equal(t, expected, ContextGetSentryCategory(ctx))
	})
}

func TestContextWithSentryCategory(t *testing.T) {
	ctx := context.Background()

	t.Run("SetCategory", func(t *testing.T) {
		expected := "category"

		ctx = ContextWithSentryCategory(ctx, expected)

		assert.Equal(t, expected, ctx.Value(sentryCategoryKey))
	})

	t.Run("OverwriteCategory", func(t *testing.T) {
		expected := "other_category"

		ctx = ContextWithSentryCategory(ctx, expected)

		assert.Equal(t, expected, ctx.Value(sentryCategoryKey))
	})
}
