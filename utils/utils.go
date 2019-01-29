package utils

import (
	"log"
	"os"
	"path/filepath"
)

// SliceStringContains sliceStringContains
func SliceStringContains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// SliceRemoveString sliceRemoveString
func SliceRemoveString(slice []string, item string) []string {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// GetSubFolders getSubFolders
func GetSubFolders(root string) []string {
	result := []string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			result = append(result, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	return result
}

// IsDirectory isDirectory
func IsDirectory(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
		return false
	}

	return fi.Mode().IsDir()
}
