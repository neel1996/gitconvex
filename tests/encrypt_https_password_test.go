package tests

import (
	"github.com/neel1996/gitconvex/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordCipherStruct_EncryptPassword(t *testing.T) {
	type fields struct {
		PlainPassword     string
		EncryptedPassword string
		KeyString         string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Tests password encryption", fields: struct {
			PlainPassword     string
			EncryptedPassword string
			KeyString         string
		}{PlainPassword: "password", EncryptedPassword: "", KeyString: "ac2f0ec3"}, want: "7XJmCQ/zn+xDqemILMJ1mS1zUqAzivF5PQ6r5YX2cEl7QhM+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := utils.PasswordCipherStruct{
				PlainPassword:     tt.fields.PlainPassword,
				EncryptedPassword: tt.fields.EncryptedPassword,
				KeyString:         tt.fields.KeyString,
			}
			if got := x.EncryptPassword(); len(got) != len(tt.want) {
				t.Errorf("EncryptPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordCipherStruct_DecryptPassword(t *testing.T) {
	type fields struct {
		PlainPassword     string
		EncryptedPassword string
		KeyString         string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Tests password decryption", fields: struct {
			PlainPassword     string
			EncryptedPassword string
			KeyString         string
		}{PlainPassword: "", EncryptedPassword: "7XJmCQ/zn+xDqemILMJ1mS1zUqAzivF5PQ6r5YX2cEl7QhM+", KeyString: "ac2f0ec3"}, want: "password"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := utils.PasswordCipherStruct{
				PlainPassword:     tt.fields.PlainPassword,
				EncryptedPassword: tt.fields.EncryptedPassword,
				KeyString:         tt.fields.KeyString,
			}
			assert.Equal(t, tt.want, x.DecryptPassword())
		})
	}
}
