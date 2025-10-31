package security

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestEncryptionManager(t *testing.T) {
	password := "test-password"

	// Test creating a new encryption manager
	t.Run("NewEncryptionManager", func(t *testing.T) {
		manager, err := NewEncryptionManager(password)
		if err != nil {
			t.Fatalf("Failed to create encryption manager: %v", err)
		}

		if manager == nil {
			t.Error("Encryption manager is nil")
		}

		if len(manager.GetSalt()) != 16 {
			t.Errorf("Expected salt length of 16, got %d", len(manager.GetSalt()))
		}
	})

	// Test encryption and decryption
	t.Run("EncryptDecrypt", func(t *testing.T) {
		manager, err := NewEncryptionManager(password)
		if err != nil {
			t.Fatalf("Failed to create encryption manager: %v", err)
		}

		plaintext := "This is a test message"

		// Encrypt the plaintext
		ciphertext, err := manager.Encrypt([]byte(plaintext))
		if err != nil {
			t.Fatalf("Failed to encrypt data: %v", err)
		}

		if len(ciphertext) <= len(plaintext) {
			t.Error("Ciphertext should be longer than plaintext due to IV and padding")
		}

		// Decrypt the ciphertext
		decrypted, err := manager.Decrypt(ciphertext)
		if err != nil {
			t.Fatalf("Failed to decrypt data: %v", err)
		}

		if string(decrypted) != plaintext {
			t.Errorf("Expected decrypted text '%s', got '%s'", plaintext, string(decrypted))
		}
	})

	// Test with existing salt
	t.Run("NewEncryptionManagerWithSalt", func(t *testing.T) {
		manager1, err := NewEncryptionManager(password)
		if err != nil {
			t.Fatalf("Failed to create encryption manager: %v", err)
		}

		salt := manager1.GetSalt()

		// Create a new manager with the same salt
		manager2 := NewEncryptionManagerWithSalt(password, salt)

		plaintext := "Test message"

		// Encrypt with first manager
		ciphertext1, err := manager1.Encrypt([]byte(plaintext))
		if err != nil {
			t.Fatalf("Failed to encrypt with manager1: %v", err)
		}

		// Decrypt with second manager
		decrypted, err := manager2.Decrypt(ciphertext1)
		if err != nil {
			t.Fatalf("Failed to decrypt with manager2: %v", err)
		}

		if string(decrypted) != plaintext {
			t.Errorf("Expected decrypted text '%s', got '%s'", plaintext, string(decrypted))
		}
	})

	// Test key file operations
	t.Run("KeyFileOperations", func(t *testing.T) {
		manager, err := NewEncryptionManager(password)
		if err != nil {
			t.Fatalf("Failed to create encryption manager: %v", err)
		}

		// Create a temporary key file
		keyFile, err := ioutil.TempFile("", "vpn_key_test")
		if err != nil {
			t.Fatalf("Failed to create temporary key file: %v", err)
		}
		defer os.Remove(keyFile.Name())
		keyFile.Close()

		// Generate key file
		if err := manager.GenerateKeyFile(keyFile.Name()); err != nil {
			t.Fatalf("Failed to generate key file: %v", err)
		}

		// Load key file
		salt, err := LoadKeyFile(keyFile.Name())
		if err != nil {
			t.Fatalf("Failed to load key file: %v", err)
		}

		if len(salt) != 16 {
			t.Errorf("Expected salt length of 16, got %d", len(salt))
		}

		// Verify the salt matches
		if string(salt) != string(manager.GetSalt()) {
			t.Error("Loaded salt does not match original salt")
		}
	})

	// Test string encryption/decryption
	t.Run("StringEncryption", func(t *testing.T) {
		manager, err := NewEncryptionManager(password)
		if err != nil {
			t.Fatalf("Failed to create encryption manager: %v", err)
		}

		plaintext := "This is a test string for encryption"

		// Encrypt string
		encrypted, err := manager.EncryptString(plaintext)
		if err != nil {
			t.Fatalf("Failed to encrypt string: %v", err)
		}

		if encrypted == plaintext {
			t.Error("Encrypted string should not be the same as plaintext")
		}

		// Decrypt string
		decrypted, err := manager.DecryptString(encrypted)
		if err != nil {
			t.Fatalf("Failed to decrypt string: %v", err)
		}

		if decrypted != plaintext {
			t.Errorf("Expected decrypted string '%s', got '%s'", plaintext, decrypted)
		}
	})
}

func TestEncryptionWithDifferentPasswords(t *testing.T) {
	password1 := "password1"
	password2 := "password2"

	manager1, err := NewEncryptionManager(password1)
	if err != nil {
		t.Fatalf("Failed to create encryption manager 1: %v", err)
	}

	manager2, err := NewEncryptionManager(password2)
	if err != nil {
		t.Fatalf("Failed to create encryption manager 2: %v", err)
	}

	plaintext := "Test data"

	// Encrypt with first manager
	ciphertext, err := manager1.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt with manager 1: %v", err)
	}

	// Try to decrypt with second manager (should fail)
	_, err = manager2.Decrypt(ciphertext)
	if err == nil {
		t.Error("Expected decryption to fail with different password")
	}
}
