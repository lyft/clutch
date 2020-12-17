package middleware

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitFullMethod(t *testing.T) {
	splitFullMethodTests := []struct {
		input   string
		service string
		method  string
		ok      bool
	}{
		{input: "", ok: false},
		{
			input:   "/foo.bar.v1.baz/FizzBuzz",
			service: "foo.bar.v1.baz",
			method:  "FizzBuzz",
			ok:      true,
		},
		{input: "baz", ok: false},
		{input: "baz/buz", ok: false},
	}

	for _, tt := range splitFullMethodTests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			s, m, ok := SplitFullMethod(tt.input)
			assert.Equal(t, tt.ok, ok)
			if tt.ok {
				assert.Equal(t, tt.service, s)
				assert.Equal(t, tt.method, m)
			} else {
				assert.Equal(t, "serviceUnknown", s)
				assert.Equal(t, "methodUnknown", m)
			}
		})
	}
}

func TestMatchFullMethod(t *testing.T) {
	input := "/foo.bar.v1.baz/FizzBuzz"
	tests := []struct {
		pattern string
		match   bool
	}{
		{pattern: "/*/*", match: true},
		{pattern: "*", match: true},
		{pattern: "", match: false},
		{pattern: "/foo.bar.v1.baz/FizzBuzz", match: true},
		{pattern: "/foo.bar.v1.baz/Fiz.Buzz", match: false},
		{pattern: "/foo.bar.v1.baz/Fizz", match: false},
		{pattern: "/foo.bar.v1.baz/fizzBuzz", match: false},
		{pattern: "/foo.bar.v1.baz/*", match: true},
		{pattern: "//*", match: false},
		{pattern: "/foo.bar.v1.baz/", match: false},
		{pattern: "/*/FizzBuzz", match: true},
		{pattern: "/blah.blah/FizzBuzz", match: false},
		{pattern: "/foo.bar.v1.baz/BlahBlah", match: false},
		{pattern: "/foo.*/FizzBuzz", match: true},
		{pattern: "/*bar*/FizzBuzz", match: true},
		{pattern: "/*/Fiz*Buzz", match: true},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			result := MatchMethodOrResource(tt.pattern, input)
			assert.Equal(t, tt.match, result)
		})
	}
}

func TestMatchResource(t *testing.T) {
	input := "us-east-1/i-12378471"
	tests := []struct {
		pattern string
		match   bool
	}{
		{pattern: "*/*", match: true},
		{pattern: "*", match: true},
		{pattern: "", match: false},
		{pattern: "us-*-1/*", match: true},
		{pattern: "us-east-1/i-12378471", match: true},
		{pattern: "us-east-1/i-12378470", match: false},
		{pattern: "*/i-12378471", match: true},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			result := MatchMethodOrResource(tt.pattern, input)
			assert.Equal(t, tt.match, result)
		})
	}
}

func TestDeeplyNestedMatching(t *testing.T) {
	input := "us-east-1/i-1234567890/network/controller"

	tests := []struct {
		pattern string
		match   bool
	}{
		{pattern: "us-east-1/**", match: true},
		{pattern: "us-east-1/i-*/**/controller", match: true},
		{pattern: "us-east-1/**/controller", match: true},
		{pattern: "us-east-1/*", match: false},
		{pattern: "us-east-1/i-*/controller", match: false},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			result := MatchMethodOrResource(tt.pattern, input)
			assert.Equal(t, tt.match, result)
		})
	}
}
