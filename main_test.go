package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCard(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  card
		err  error
	}{
		{
			name: "Parses a card with quantity 1 without error",
			in:   "1 Swooping Lookout",
			out: card{
				name:     "Swooping Lookout",
				quantity: 1,
			},
			err: nil,
		},
		{
			name: "Parses a card with quantity > 1 without error",
			in:   "5 Swooping Lookout",
			out: card{
				name:     "Swooping Lookout",
				quantity: 5,
			},
			err: nil,
		},
		{
			name: "Parses a card with expansion information without error",
			in:   "1 Terramorphic Expanse (ONE) 261			",
			out: card{
				name:     "Terramorphic Expanse",
				quantity: 1,
			},
			err: nil,
		},
		{
			name: "errors on invalid card format",
			in:   "somethingWithoutASpace",
			out:  card{},
			err:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := parseCard(tc.in)
			if tc.err == nil {
				assert.Equal(t, tc.err, err)
			} else {
				assert.Equal(t, tc.out, out)
			}

		})

	}
}
