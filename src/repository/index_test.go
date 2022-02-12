package repository

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/myrachanto/ddd/httperors"
)

func TestGetconnected(t *testing.T) {
	GormDB, err := IndexRepo.Getconnected()
	assert.NotNil(t, GormDB, "The conncetion passed")
	assert.Nil(t, err, "The error must be nil for successifully database connection")
}
