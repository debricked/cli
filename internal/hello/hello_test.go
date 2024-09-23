package hello

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDebrickedGreeter(t *testing.T) {
	greeter := NewDebrickedGreeter()
	assert.NotNil(t, greeter)
}

func TestGreeting(t *testing.T) {
	greeter := NewDebrickedGreeter()
	greeting := greeter.Greeting("Debricked")
	assert.Equal(t, greeting, "Hello Debricked!")
}
