package models

import "testing"

func TestEncryptAndDecrypt(t *testing.T) {
	// Define test data
	var a, b, c, d = &Login{Info: "test_username"},
		&BankCard{Info: "test_card_number"},
		&ArbitraryText{Text: "test_text"},
		&Binary{Info: "test_info"}

	// Define test passphrase
	passphrase := "my passphrase"

	// Encrypt the data
	encryptedA, err := a.EncryptAll(passphrase)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	encryptedB, err := b.EncryptAll(passphrase)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	encryptedC, err := c.EncryptAll(passphrase)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	encryptedD, err := d.EncryptAll(passphrase)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Decrypt the data
	decryptedA := &Login{}
	err = decryptedA.DecryptAll(passphrase, encryptedA)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	decryptedB := &BankCard{}
	err = decryptedB.DecryptAll(passphrase, encryptedB)
	if err != nil {

		t.Errorf("Unexpected error: %v", err)
	}

	decryptedC := &ArbitraryText{}
	err = decryptedC.DecryptAll(passphrase, encryptedC)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	decryptedD := &Binary{}
	err = decryptedD.DecryptAll(passphrase, encryptedD)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the decrypted data matches the original data

	if decryptedA.Info != a.Info {
		t.Errorf("Decrypted data does not match original data")
	}

	if decryptedB.Info != b.Info {

		t.Errorf("Decrypted data does not match original data")
	}

	if decryptedC.Text != c.Text {

		t.Errorf("Decrypted data does not match original data")
	}

	if decryptedD.Info != d.Info {

		t.Errorf("Decrypted data does not match original data")
	}

}
