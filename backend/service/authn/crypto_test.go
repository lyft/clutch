package authn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadKeyLength(t *testing.T) {
	c, err := newCryptographer("")
	assert.Nil(t, c)
	assert.Error(t, err)
}

func TestEmpty(t *testing.T) {
	c, _ := newCryptographer("hello world")
	_, err := c.Encrypt([]byte{})
	assert.Error(t, err)
	_, err = c.Encrypt(nil)
	assert.Error(t, err)
}

func TestKeyMakesADifference(t *testing.T) {
	in := []byte("zoom")

	c1, _ := newCryptographer("abcdefgh")
	o1, _ := c1.Encrypt(in)

	c2, _ := newCryptographer("abcdefghi")
	o2, _ := c2.Encrypt(in)

	assert.NotEqual(t, o1, o2)
}

func TestRoundTrip(t *testing.T) {
	c, err := newCryptographer("123456")
	assert.NoError(t, err)

	tests := []struct{ in string }{
		{in: "foo"},
		{in: "bar"},
		{in: "foobar"},
		{in: "foo bar"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			b := []byte(tt.in)
			encrypted, err := c.Encrypt(b)
			assert.NoError(t, err)

			decrypted, err := c.Decrypt(encrypted)
			assert.NoError(t, err)
			assert.Equal(t, b, decrypted)
		})
	}
}
