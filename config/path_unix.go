// +build !windows

package config

import (
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const (
	defaultFileName    = ".config/prm"
	oldDefaultFileName = ".prm"
)

// GetPath get the configuration file path.
func GetPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	pathOne := filepath.Join(usr.HomeDir, oldDefaultFileName)
	pathTwo := filepath.Join(usr.HomeDir, defaultFileName)

	info, err := os.Stat(pathOne)
	if err != nil {
		if os.IsNotExist(err) {
			return pathTwo, nil
		}
		return "", err
	}

	log.Println("WARN: old configuration file detected, migration in progress.")

	err = copyConfigFile(pathOne, pathTwo, info)
	if err != nil {
		return "", err
	}

	err = os.Remove(pathOne)
	if err != nil {
		return "", err
	}

	log.Println("WARN: old configuration file detected, migration done.")

	return pathTwo, nil
}

func copyConfigFile(src, dst string, info os.FileInfo) error {
	baseDir := path.Dir(dst)
	err := os.MkdirAll(baseDir, 0700)
	if err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer safeClose(f.Close)

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer safeClose(s.Close)

	_, err = io.Copy(f, s)
	return err
}

func safeClose(fn func() error) {
	err := fn()
	if err != nil {
		log.Println(err)
	}
}
