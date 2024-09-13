package keyring

import (
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKeySet(t *testing.T) {
	type args struct {
		keyLifespan int
		jwtLifespan int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid KeySet",
			args: args{
				keyLifespan: 129600,
				jwtLifespan: 86400,
			},
			wantErr: false,
		},
		{
			name: "Invalid KeySet",
			args: args{
				keyLifespan: 86400,
				jwtLifespan: 129600,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSet(tt.args.keyLifespan, tt.args.jwtLifespan, hclog.Default())
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSet() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestKeySet_SetCurveTypes(t *testing.T) {
	type args struct {
		curveTypes []CurveType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should Set Generation Types",
			args: args{
				curveTypes: []CurveType{P256, P384, P521},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyset, err := NewSet(200, 100, hclog.Default())
			if err != nil {
				t.Errorf("NewKeySet() error = %v", err)
				return
			}

			keyset.SetCurveTypes(tt.args.curveTypes...)
			assert.Equal(t, tt.args.curveTypes, keyset.curveTypes)
		})
	}
}

func TestKeySet_Generate(t *testing.T) {
	type args struct {
		curveType CurveType
		n         int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should generate P-256",
			args: args{
				curveType: P256,
				n:         3,
			},
		},
		{
			name: "Should generate P-384",
			args: args{
				curveType: P384,
				n:         3,
			},
		},
		{
			name: "Should generate P-521",
			args: args{
				curveType: P521,
				n:         3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyset, _ := NewSet(1000, 500, hclog.Default())
			keyset.SetCurveTypes(tt.args.curveType)

			err := keyset.Generate(tt.args.n)
			if err != nil {
				t.Errorf("Generate() error = %v", err)
				return
			}

			for _, key := range keyset.GetPrivateKeySet().Keys {
				assert.NotNil(t, keyset.activeKeys[key.KeyID()], "should be active")
				assert.NotNil(t, keyset.expiringKeys[key.KeyID()], "should be expiring")
				assert.Equal(t, tt.args.curveType.GetAlgorithm(), key.Algorithm())
			}
		})
	}
}
