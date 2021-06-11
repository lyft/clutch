package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertLockToUint32(t *testing.T) {
	tests := []struct {
		id     string
		input  string
		expect uint32
	}{
		{
			id:     "key with chars",
			input:  "audit:eventsink",
			expect: 1635083369,
		},
		{
			id:     "key with special chars",
			input:  "*()#@&!*(#!@",
			expect: 707275043,
		},
		{
			id:     "key with numbers",
			input:  "123123123",
			expect: 825373489,
		},
	}

	for _, tt := range tests {
		id := ConvertLockToUint32(tt.input)
		assert.Equal(t, tt.expect, id)
	}
}
