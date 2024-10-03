package bpti_helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DirectoryHasPages(directory_path, file_extension string) (bool, error) {
	var has_pngs bool = false
	directory_entries, err := os.ReadDir(directory_path)
	if err != nil {
		return has_pngs, err
	}

	for _, entry := range directory_entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), file_extension) {
			has_pngs = true
			break
		}
	}

	return has_pngs, nil
}

func CleanDirectoryPreviousPages(directory_path, file_extension string) error {
	directory_entries, err := os.ReadDir(directory_path)
	if err != nil {
		return err
	}

	for _, entry := range directory_entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), file_extension) {
			err = os.Remove(filepath.Join(directory_path, entry.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateDirectoryIfNotExists(directory_path string) error {
	if PathExists(directory_path) {
		return nil
	}

	err := os.Mkdir(directory_path, 0755)
	if err != nil {
		return err
	}

	return nil
}

func GetParentDirectory(path string) string {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	parent_dir, _ := filepath.Split(path)

	return parent_dir
}

func GetPDFDirectoryName(pdf_path string) string {
	var pdf_name, pdf_ext string

	_, pdf_name = filepath.Split(pdf_path)
	pdf_ext = filepath.Ext(pdf_name)

	directory_base_name := strings.TrimSuffix(pdf_name, pdf_ext)

	return directory_base_name
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDirectory(path string) bool {
	file_info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return file_info.IsDir()
}

func IsPDF(path string) (bool, error) {
	if path_exists := PathExists(path); !path_exists {
		return false, fmt.Errorf("Path does not exist: %s", path)
	}

	ext := filepath.Ext(path)

	return ext == ".pdf", nil
}
