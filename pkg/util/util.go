package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// VerifyOrCreateFile checks if a file exists at the given path,
// and if not, creates all parent directories and creates the file.
func VerifyOrCreateFile(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if err = os.MkdirAll(dir, 0o744); err != nil {
			return fmt.Errorf("creating parent directory: %w", err)
		}
		if _, err = os.Create(path); err != nil {
			return fmt.Errorf("creating file: %w", err)
		}
		return nil
	}
	return err
}
