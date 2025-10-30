package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"golang.org/x/crypto/pbkdf2"
)

// EncryptionManager manages data encryption and decryption
type EncryptionManager struct {
	key    []byte
	salt   []byte
}

// NewEncryptionManager creates a new encryption manager
func NewEncryptionManager(password string) (*EncryptionManager, error) {
	// Generate a salt
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from password using PBKDF2
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)

	return &EncryptionManager{
		key:  key,
		salt: salt,
	}, nil
}

// NewEncryptionManagerWithSalt creates a new encryption manager with existing salt
func NewEncryptionManagerWithSalt(password string, salt []byte) *EncryptionManager {
	// Derive key from password using PBKDF2
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)

	return &EncryptionManager{
		key:  key,
		salt: salt,
	}
}

// GetSalt returns the salt used for key derivation
func (em *EncryptionManager) GetSalt() []byte {
	return em.salt
}

// Encrypt encrypts plaintext data
func (em *EncryptionManager) Encrypt(plaintext []byte) ([]byte, error) {
	// Create a new AES cipher
	block, err := aes.NewCipher(em.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	// Create a new CBC encrypter
	mode := cipher.NewCBCEncrypter(block, iv)

	// Pad the plaintext to be a multiple of the block size
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	paddedPlaintext := make([]byte, len(plaintext)+padding)
	copy(paddedPlaintext, plaintext)
	for i := len(plaintext); i < len(paddedPlaintext); i++ {
		paddedPlaintext[i] = byte(padding)
	}

	// Encrypt the padded plaintext
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// Prepend the IV to the ciphertext
	result := make([]byte, len(iv)+len(ciphertext))
	copy(result, iv)
	copy(result[len(iv):], ciphertext)

	return result, nil
}

// Decrypt decrypts ciphertext data
func (em *EncryptionManager) Decrypt(ciphertext []byte) ([]byte, error) {
	// Create a new AES cipher
	block, err := aes.NewCipher(em.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Extract the IV from the ciphertext
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Create a new CBC decrypter
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the ciphertext
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove padding
	padding := int(plaintext[len(plaintext)-1])
	if padding > len(plaintext) {
		return nil, fmt.Errorf("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-padding]

	return plaintext, nil
}

// EncryptFile encrypts a file
func (em *EncryptionManager) EncryptFile(inputPath, outputPath string) error {
	// Read the input file
	plaintext, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Encrypt the data
	ciphertext, err := em.Encrypt(plaintext)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Write the encrypted data to the output file
	if err := ioutil.WriteFile(outputPath, ciphertext, 0600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// DecryptFile decrypts a file
func (em *EncryptionManager) DecryptFile(inputPath, outputPath string) error {
	// Read the input file
	ciphertext, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Decrypt the data
	plaintext, err := em.Decrypt(ciphertext)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Write the decrypted data to the output file
	if err := ioutil.WriteFile(outputPath, plaintext, 0600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// EncryptDatabase encrypts the database file
func (em *EncryptionManager) EncryptDatabase(dbPath, encryptedPath string) error {
	return em.EncryptFile(dbPath, encryptedPath)
}

// DecryptDatabase decrypts the database file
func (em *EncryptionManager) DecryptDatabase(encryptedPath, dbPath string) error {
	return em.DecryptFile(encryptedPath, dbPath)
}

// GenerateKeyFile generates a key file with the salt
func (em *EncryptionManager) GenerateKeyFile(keyPath string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(keyPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the salt to the key file
	if err := ioutil.WriteFile(keyPath, em.salt, 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

// LoadKeyFile loads the salt from a key file
func LoadKeyFile(keyPath string) ([]byte, error) {
	// Read the salt from the key file
	salt, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	return salt, nil
}

// EncryptString encrypts a string and returns a base64 encoded string
func (em *EncryptionManager) EncryptString(plaintext string) (string, error) {
	ciphertext, err := em.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptString decrypts a base64 encoded string
func (em *EncryptionManager) DecryptString(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	
	plaintext, err := em.Decrypt(data)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}