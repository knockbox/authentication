package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCryptoRandom(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should generate random in range",
			args: args{
				min: 1,
				max: 1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := tt.args.min; i < tt.args.max; i++ {
				n, err := CryptoRandom(i)
				if err != nil {
					t.Errorf("CryptoRandom() error = %v", err)
					return
				}

				assert.True(t, int(n) < i, "rand should be in range [0, i)")
			}
		})
	}
}
