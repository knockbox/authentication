package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPayloadHasChanges(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should have changes",
			args: args{
				payload: struct {
					field1 *string
					field2 *int
					field3 *bool
				}{
					field1: new(string),
					field2: new(int),
					field3: new(bool),
				},
			},
			want: true,
		},
		{
			name: "should have no changes",
			args: args{
				payload: struct {
					field1 *string
					field2 *int
					field3 *bool
				}{
					field1: nil,
					field2: nil,
					field3: nil,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, PayloadHasChanges(tt.args.payload), "PayloadHasChanges(%v)", tt.args.payload)
		})
	}
}
