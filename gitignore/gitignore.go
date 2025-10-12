package gitignore

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ensureFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			return fmt.Errorf("failed to create %v: %w", path, err)
		}
	}
	return nil
}

func lineExistsInFile(path, target string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == target {
			return true, nil
		}
	}
	return false, scanner.Err()
}

func AddToGitignore(rootPath, target string) error {
	gitignorePath := filepath.Join(rootPath, ".gitignore")
	err := ensureFileExists(gitignorePath)
	if err != nil {
		return err
	}
	lineExistsInFile, err := lineExistsInFile(gitignorePath, target)
	if lineExistsInFile {
		return nil
	}
	if err != nil {
		return err
	}

	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(target + "\n"); err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
	}

	fmt.Printf("🆕 Added '%s' to .gitignore\n", target)
	return nil
}
