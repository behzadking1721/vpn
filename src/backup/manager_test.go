package backup

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/security"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBackupManager(t *testing.T) {
	// Create temporary directories for testing
	tempDir, err := ioutil.TempDir("", "vpn_backup_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "vpn.db")
	backupDir := filepath.Join(tempDir, "backups")

	// Create a test database file
	dbContent := "This is test database content"
	if err := ioutil.WriteFile(dbPath, []byte(dbContent), 0644); err != nil {
		t.Fatalf("Failed to create test database file: %v", err)
	}

	// Create encryption manager
	encryptManager, err := security.NewEncryptionManager("test-password")
	if err != nil {
		t.Fatalf("Failed to create encryption manager: %v", err)
	}

	// Create backup manager
	backupManager := NewBackupManager(dbPath, backupDir, encryptManager)

	// Test creating a backup without encryption
	t.Run("CreateBackupWithoutEncryption", func(t *testing.T) {
		backup, err := backupManager.CreateBackup(false)
		if err != nil {
			t.Fatalf("Failed to create backup: %v", err)
		}

		if backup == nil {
			t.Error("Backup is nil")
		}

		if backup.Name == "" {
			t.Error("Backup name is empty")
		}

		if backup.Size <= 0 {
			t.Error("Backup size is not positive")
		}

		if backup.Checksum == "" {
			t.Error("Backup checksum is empty")
		}

		if backup.Encrypted {
			t.Error("Backup should not be marked as encrypted")
		}
	})

	// Test creating a backup with encryption
	t.Run("CreateBackupWithEncryption", func(t *testing.T) {
		backup, err := backupManager.CreateBackup(true)
		if err != nil {
			t.Fatalf("Failed to create encrypted backup: %v", err)
		}

		if backup == nil {
			t.Error("Encrypted backup is nil")
		}

		if backup.Encrypted == false {
			t.Error("Backup should be marked as encrypted")
		}
	})

	// Test listing backups
	t.Run("ListBackups", func(t *testing.T) {
		backups, err := backupManager.ListBackups()
		if err != nil {
			t.Fatalf("Failed to list backups: %v", err)
		}

		if len(backups) < 2 {
			t.Errorf("Expected at least 2 backups, got %d", len(backups))
		}

		// Check that all backups have the required fields
		for _, backup := range backups {
			if backup.ID == "" {
				t.Error("Backup ID is empty")
			}

			if backup.Name == "" {
				t.Error("Backup name is empty")
			}

			if backup.Timestamp.IsZero() {
				t.Error("Backup timestamp is zero")
			}

			if backup.Size <= 0 {
				t.Error("Backup size is not positive")
			}

			if backup.Checksum == "" {
				t.Error("Backup checksum is empty")
			}
		}
	})

	// Test deleting a backup
	t.Run("DeleteBackup", func(t *testing.T) {
		// Create a backup to delete
		backup, err := backupManager.CreateBackup(false)
		if err != nil {
			t.Fatalf("Failed to create backup for deletion: %v", err)
		}

		// Delete the backup
		err = backupManager.DeleteBackup(backup.Name)
		if err != nil {
			t.Fatalf("Failed to delete backup: %v", err)
		}

		// Try to delete the same backup again (should fail)
		err = backupManager.DeleteBackup(backup.Name)
		if err == nil {
			t.Error("Expected error when deleting non-existent backup")
		}
	})
}

func TestBackupRestore(t *testing.T) {
	// Create temporary directories for testing
	tempDir, err := ioutil.TempDir("", "vpn_restore_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "vpn.db")
	backupDir := filepath.Join(tempDir, "backups")

	// Create a test database file
	originalContent := "Original database content"
	if err := ioutil.WriteFile(dbPath, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create test database file: %v", err)
	}

	// Create backup manager without encryption for simplicity
	backupManager := NewBackupManager(dbPath, backupDir, nil)

	// Create a backup
	backup, err := backupManager.CreateBackup(false)
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	// Modify the database file
	modifiedContent := "Modified database content"
	if err := ioutil.WriteFile(dbPath, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to modify database file: %v", err)
	}

	// Verify the content was modified
	content, err := ioutil.ReadFile(dbPath)
	if err != nil {
		t.Fatalf("Failed to read database file: %v", err)
	}

	if string(content) != modifiedContent {
		t.Errorf("Expected modified content '%s', got '%s'", modifiedContent, string(content))
	}

	// Restore the backup
	err = backupManager.RestoreBackup(backup.Name)
	if err != nil {
		t.Fatalf("Failed to restore backup: %v", err)
	}

	// Verify the content was restored
	restoredContent, err := ioutil.ReadFile(dbPath)
	if err != nil {
		t.Fatalf("Failed to read restored database file: %v", err)
	}

	if string(restoredContent) != originalContent {
		t.Errorf("Expected restored content '%s', got '%s'", originalContent, string(restoredContent))
	}
}

func TestBackupManagerWithoutEncryption(t *testing.T) {
	// Create temporary directories for testing
	tempDir, err := ioutil.TempDir("", "vpn_backup_no_encrypt_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "vpn.db")
	backupDir := filepath.Join(tempDir, "backups")

	// Create a test database file
	dbContent := "Test database content without encryption"
	if err := ioutil.WriteFile(dbPath, []byte(dbContent), 0644); err != nil {
		t.Fatalf("Failed to create test database file: %v", err)
	}

	// Create backup manager without encryption
	backupManager := NewBackupManager(dbPath, backupDir, nil)

	// Test creating a backup without encryption
	backup, err := backupManager.CreateBackup(false)
	if err != nil {
		t.Fatalf("Failed to create backup without encryption: %v", err)
	}

	if backup.Encrypted {
		t.Error("Backup should not be marked as encrypted when encryption is disabled")
	}

	// Test creating a backup with encryption when encryption manager is nil
	backup, err = backupManager.CreateBackup(true)
	if err != nil {
		t.Fatalf("Failed to create backup with encryption disabled: %v", err)
	}

	if backup.Encrypted {
		t.Error("Backup should not be marked as encrypted when encryption manager is nil")
	}
}
