package backup

import (
	"archive/zip"
	"c:/Users/behza/OneDrive/Documents/vpn/src/security"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// BackupManager manages database backup and restore operations
type BackupManager struct {
	dbPath     string
	backupDir  string
	encryptMgr *security.EncryptionManager
}

// BackupInfo contains information about a backup
type BackupInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Size      int64     `json:"size"`
	Checksum  string    `json:"checksum"`
	Encrypted bool      `json:"encrypted"`
}

// NewBackupManager creates a new backup manager
func NewBackupManager(dbPath, backupDir string, encryptMgr *security.EncryptionManager) *BackupManager {
	return &BackupManager{
		dbPath:     dbPath,
		backupDir:  backupDir,
		encryptMgr: encryptMgr,
	}
}

// CreateBackup creates a new backup of the database
func (bm *BackupManager) CreateBackup(encrypt bool) (*BackupInfo, error) {
	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(bm.backupDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename
	timestamp := time.Now()
	filename := fmt.Sprintf("vpn_backup_%s.zip", timestamp.Format("20060102_150405"))
	backupPath := filepath.Join(bm.backupDir, filename)

	// Create zip file
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer zipFile.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add database file to zip
	dbFile, err := os.Open(bm.dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database file: %w", err)
	}
	defer dbFile.Close()

	// Get database file info
	dbInfo, err := dbFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get database file info: %w", err)
	}

	// Create file header in zip
	dbHeader, err := zip.FileInfoHeader(dbInfo, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create file header: %w", err)
	}
	dbHeader.Name = "vpn.db"

	// Create writer for the file in zip
	dbWriter, err := zipWriter.CreateHeader(dbHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to create writer for database file: %w", err)
	}

	// Copy database file to zip
	var reader io.Reader = dbFile
	if encrypt && bm.encryptMgr != nil {
		// If encryption is enabled, encrypt the database file
		dbData, err := ioutil.ReadAll(dbFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read database file: %w", err)
		}

		encryptedData, err := bm.encryptMgr.Encrypt(dbData)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt database: %w", err)
		}

		reader = &byteReader{data: encryptedData}
	}

	if _, err := io.Copy(dbWriter, reader); err != nil {
		return nil, fmt.Errorf("failed to copy database to backup: %w", err)
	}

	// Close zip writer to finalize the zip file
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}

	// Get backup file info
	backupInfo, err := os.Stat(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup file info: %w", err)
	}

	// Calculate checksum
	checksum, err := calculateChecksum(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// Create backup info
	backup := &BackupInfo{
		ID:        timestamp.Format("20060102150405"),
		Name:      filename,
		Timestamp: timestamp,
		Size:      backupInfo.Size(),
		Checksum:  checksum,
		Encrypted: encrypt && bm.encryptMgr != nil,
	}

	return backup, nil
}

// RestoreBackup restores a backup
func (bm *BackupManager) RestoreBackup(backupName string) error {
	// Open backup file
	backupPath := filepath.Join(bm.backupDir, backupName)
	zipReader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer zipReader.Close()

	// Find database file in backup
	var dbFile *zip.File
	for _, file := range zipReader.File {
		if file.Name == "vpn.db" {
			dbFile = file
			break
		}
	}

	if dbFile == nil {
		return fmt.Errorf("database file not found in backup")
	}

	// Open database file in backup
	dbReader, err := dbFile.Open()
	if err != nil {
		return fmt.Errorf("failed to open database file in backup: %w", err)
	}
	defer dbReader.Close()

	// Create temporary file for restored database
	tempPath := bm.dbPath + ".restore"
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Copy data from backup to temporary file
	var reader io.Reader = dbReader
	if bm.encryptMgr != nil && dbFile.FileInfo().Size() > 0 {
		// If encryption manager exists, try to decrypt
		dbData, err := ioutil.ReadAll(dbReader)
		if err != nil {
			return fmt.Errorf("failed to read database from backup: %w", err)
		}

		decryptedData, err := bm.encryptMgr.Decrypt(dbData)
		if err != nil {
			// If decryption fails, assume the backup is not encrypted
			// and write the original data
			reader = &byteReader{data: dbData}
		} else {
			reader = &byteReader{data: decryptedData}
		}
	}

	if _, err := io.Copy(tempFile, reader); err != nil {
		return fmt.Errorf("failed to copy data to temporary file: %w", err)
	}

	// Close temporary file
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Replace original database file with restored one
	if err := os.Rename(tempPath, bm.dbPath); err != nil {
		// If rename fails, try to copy
		if err := copyFile(tempPath, bm.dbPath); err != nil {
			return fmt.Errorf("failed to replace database file: %w", err)
		}

		// Remove temporary file
		os.Remove(tempPath)
	}

	return nil
}

// ListBackups lists all available backups
func (bm *BackupManager) ListBackups() ([]*BackupInfo, error) {
	// Check if backup directory exists
	if _, err := os.Stat(bm.backupDir); os.IsNotExist(err) {
		return []*BackupInfo{}, nil
	}

	// Read backup directory
	files, err := ioutil.ReadDir(bm.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	// Create list of backups
	backups := make([]*BackupInfo, 0)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".zip" {
			// Parse timestamp from filename
			// Expected format: vpn_backup_20060102_150405.zip
			timestampStr := file.Name()[11 : len(file.Name())-4] // Remove "vpn_backup_" and ".zip"
			timestamp, err := time.Parse("20060102_150405", timestampStr)
			if err != nil {
				// Skip files with invalid format
				continue
			}

			// Calculate checksum
			checksum, err := calculateChecksum(filepath.Join(bm.backupDir, file.Name()))
			if err != nil {
				// Skip files with checksum errors
				continue
			}

			// Create backup info
			backup := &BackupInfo{
				ID:        timestamp.Format("20060102150405"),
				Name:      file.Name(),
				Timestamp: timestamp,
				Size:      file.Size(),
				Checksum:  checksum,
				Encrypted: false, // We can't determine this without examining the file content
			}

			backups = append(backups, backup)
		}
	}

	return backups, nil
}

// DeleteBackup deletes a backup
func (bm *BackupManager) DeleteBackup(backupName string) error {
	backupPath := filepath.Join(bm.backupDir, backupName)

	// Check if file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	// Delete the file
	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}

	return nil
}

// calculateChecksum calculates the SHA256 checksum of a file
func calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// byteReader is a simple io.Reader implementation for byte slices
type byteReader struct {
	data []byte
	pos  int
}

func (br *byteReader) Read(p []byte) (n int, err error) {
	if br.pos >= len(br.data) {
		return 0, io.EOF
	}

	n = copy(p, br.data[br.pos:])
	br.pos += n

	if br.pos >= len(br.data) {
		err = io.EOF
	}

	return n, err
}
