package main

import (
	"testing"

	"fyne.io/fyne/test"

	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
	"github.com/stretchr/testify/assert"
)

func TestCompose_Send(t *testing.T) {
	server := client.NewTestServer()
	compose := newCompose(test.NewApp(), server)
	ui := compose.loadUI()

	pending := len(server.Outbox)
	closed := false
	ui.SetOnClosed(func() {
		closed = true
	})
	address := "test@example.com"
	test.Type(compose.to, address)

	test.Click(compose.send)
	assert.True(t, closed)
	assert.Equal(t, pending + 1, len(server.Outbox))

	email := server.Outbox[len(server.Outbox)-1]
	assert.Equal(t, address, email.ToEmailString())
}
