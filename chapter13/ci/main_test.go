package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_RunApp(t *testing.T) {
	testStore := newTestStorage()
	runApp(testStore)

	newFile := testStore.Read("README.txt")
	assert.Equal(t, "overwrite", newFile)
}
