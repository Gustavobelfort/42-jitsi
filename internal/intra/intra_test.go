package intra

import (
	"testing"

	"github.com/gustavobelfort/42-jitsi/internal/config"
	"github.com/magiconair/properties/assert"
)

func TestGetToken(t *testing.T) {

	config.Initiate()

	client, err := New()
	if err != nil {
		t.Errorf("Failed to initialize the client with error: %d", err)
	}

	err = client.GetToken()
	if err != nil {
		t.Errorf("Failed to Get the Token with error: %+x", err.Error())
	}
}

func TestGetUserEmail(t *testing.T) {

	config.Initiate()

	client, err := New()
	client.GetToken()

	if err != nil {
		t.Errorf("Failed to initialize the client with error: %d", err)
	}

	email, err := client.GetUserEmail("gbelfort")
	if err != nil {
		t.Errorf("Failed to Get the user email with error: %s", err.Error())
	}

	assert.Equal(t, email, "gbelfort@student.42.us.org")
}
