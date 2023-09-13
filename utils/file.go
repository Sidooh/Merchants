package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func GetFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Error(err)
		return nil
	}

	return file
}

func GetLogFile(filename string) *os.File {
	pwd, err := os.Getwd()

	file := GetFile(filepath.Join(pwd, "logs/", filename))
	if err != nil || file == nil {
		file = os.Stdout
	}

	return file
}
