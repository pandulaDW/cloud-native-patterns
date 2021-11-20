package testing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddTable(t *testing.T) {
	cases := []struct {
		Name           string
		A, B, Expected int
	}{
		{"foo", 1, 1, 2},
		{"boo", 1, -1, 0},
		{"bar", 1, 0, 1},
		{"baz", 0, 0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, tc.A+tc.B)
		})
	}
}
