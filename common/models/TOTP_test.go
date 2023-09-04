package models

import (
	"github.com/pquerna/otp/totp"
	"testing"
)

func TestGenerateOneTimePassword(t *testing.T) {
	// Gen random key
	generate, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "goph-keeper",
		AccountName: "testuser",
	})
	if err != nil {
		return
	}

	// Create a Login instance
	login := &Login{
		Username: "testuser",
	}
	ok := login.RegisterOneTime(generate.Secret())
	if !ok {
		return
	}
	// Generate a one-time password
	otp, _, _ := login.GenerateOneTimePassword()
	valid := totp.Validate(otp, generate.Secret())
	if !valid {
		t.Error("Error validating OTP")
	}
}
