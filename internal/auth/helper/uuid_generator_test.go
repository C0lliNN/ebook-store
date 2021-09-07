// +build unit

package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUIDGenerator_NewID(t *testing.T) {
	id := UUIDGenerator{}.NewID()

	assert.Len(t, id, 36)
}
