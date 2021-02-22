package helper

import (
	"os"
	"path/filepath"
)

func chDir(path string) error {
	os.Chdir(path)
	_, err := os.Getwd()
	if err != nil {
		return err
	}

	return nil
}

func walkMatch(start string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(start,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				dirs = append(dirs, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}
