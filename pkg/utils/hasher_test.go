package utils

import (
	"reflect"
	"testing"
)

func TestComparePasswords(t *testing.T) {
	type args struct {
		hashedPassword   string
		unhashedPassword string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should compare passwords true",
			args: args{
				hashedPassword:   "$2a$12$ih4hpok3suHDbU5ob26LDead1OtqDKlau7XinAzTIvNDz8r7Vi7zC",
				unhashedPassword: "hello",
			},
			want: true,
		},
		{
			name: "should compare passwords false",
			args: args{
				hashedPassword:   "$2a$12$ih4hpok3suHDbU5ob26LDead1OtqDKlau7XinAzTIvNDz8r7Vi7zC",
				unhashedPassword: "heyo",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComparePasswords(tt.args.hashedPassword, tt.args.unhashedPassword); got != tt.want {
				t.Errorf("ComparePasswords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should generate password",
			args: args{
				password: "password-123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GeneratePassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNormalizePassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "should normalize password",
			args: args{password: "password-123"},
			want: []byte("password-123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePassword(tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalizePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
