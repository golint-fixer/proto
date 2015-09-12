package proto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCopy(t *testing.T) {
	assert := assert.New(t)

	var a []byte

	b := []byte{0, 1, 2, 3}

	c := Copy(a, b)

	E := []byte{0, 1, 2, 3}

	assert.Equal(
		E,
		c,
		"bytes are equal",
	)

	c = Copy(c, a)

	E = []byte{}

	assert.Equal(
		E,
		c,
		"there are no bytes",
	)
}
