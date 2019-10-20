package main_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	_, err := os.Stat("/opt/resource/check")
	assert.NoError(t, err, "Compiled 'check' binary doesn't exist")
}
