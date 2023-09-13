package models

import (
	"testing"

	"github.com/pquerna/otp/totp"
)

func TestGenerateOneTimePassword(t *testing.T) {
	// Gen random key
	generate, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "goph-keeper",
		AccountName: "testuser",
	})
	if err != nil {
		t.Error("Error generating OTP")
	}

	// Create a Login instance
	login := &Login{
		Username: "testuser",
	}
	ok := login.RegisterOneTime(generate.Secret())
	if !ok {
		t.Error("Error registering one-time password")
	}
	// Generate a one-time password
	otp, _, _ := login.GenerateOneTimePassword()
	valid := totp.Validate(otp, generate.Secret())
	if !valid {
		t.Error("Error validating OTP")
	}
}

func TestRegisterOneTime(t *testing.T) {
	// Create a new login instance
	login := &Login{}

	// Generate a secret
	generate, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "goph-keeper",
		AccountName: "testuser",
	})
	if err != nil {
		t.Error("Error generating OTP")
	}

	// Register a new one-time password
	success := login.RegisterOneTime(generate.Secret())
	if !success {
		t.Errorf("Failed to register one-time password")
	}

	// Check if the one-time origin matches the secret
	if login.OneTimeOrigin != generate.Secret() {
		t.Errorf("One-time origin does not match secret")
	}
}

func TestRegisterRecoveryCodes(t *testing.T) {
	// Create a new login instance
	login := &Login{}

	// Define test recovery codes
	recoveryCodes := "test recovery codes"

	// Register the recovery codes
	login.RegisterRecoveryCodes(recoveryCodes)

	// Check if the recovery codes match the test recovery codes
	if login.RecoveryCodes != recoveryCodes {
		t.Errorf("Recovery codes do not match test recovery codes")
	}
}
