package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Random capitalisation
	headers = NewHeaders()
	data = []byte("HoSt: localhost:8080\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:8080", headers["host"])
	assert.Equal(t, 22, n)
	assert.False(t, done)

	// Test: Invalid chracter header
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Many Keys
	headers = NewHeaders()
	data = []byte("Host: localhost:8080\r\nHost:localhost:6969\r\nHost: localhost:6767\r\n\r\n")
	n1, done, err := headers.Parse(data)
	n2, done, err := headers.Parse(data[n1:])
	n3, done, err := headers.Parse(data[n1+n2:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:8080, localhost:6969, localhost:6767", headers["host"])
	assert.Equal(t, 65, n1+n2+n3)
	assert.False(t, done)
}
