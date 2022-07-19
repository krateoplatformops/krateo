package osutils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	if info.IsDir() {
		return false, fmt.Errorf("%v: is a directory, expected file", path)
	}

	return true, nil
}

func GetAppDir(appName string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, fmt.Sprintf(".%s", appName))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return dir, nil
}
