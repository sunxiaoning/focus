package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
)

type fileHelper struct {
}

var FileHelper = &fileHelper{}

type fileOperation interface {
	CreateDirectory(dirPath string) (err error)
	OpenFile(filePath string, flag int, perm os.FileMode) (file *os.File, err error)
	PathExist(filePath string) (exists bool)
}

func (fileHelper *fileHelper) CreateDirectory(dirPath string) (err error) {
	f, e := os.Stat(dirPath)
	if e != nil && os.IsNotExist(e) {
		return os.MkdirAll(dirPath, 0755)
	}
	if e == nil && !f.IsDir() {
		return fmt.Errorf("create dir:%s error, not a directory", dirPath)
	}
	return e
}
func (fileHelper *fileHelper) OpenFile(filePath string, flag int, perm os.FileMode) (file *os.File, err error) {
	if fileHelper.PathExist(filePath) {
		return os.OpenFile(filePath, flag, perm)
	}
	if err := fileHelper.CreateDirectory(filepath.Dir(filePath)); err != nil {
		return nil, err
	}

	return os.OpenFile(filePath, flag, perm)
}

func (fileHelper *fileHelper) DefOpenFile(filePath string) (file *os.File, err error) {
	return fileHelper.OpenFile(filePath, os.O_RDONLY, 0)
}

func (fileHelper *fileHelper) PathExist(filePath string) (exists bool) {
	_, err := os.Stat(filePath)
	return err == nil
}

func CreateDirectory(dirPath string) (err error) {
	return FileHelper.CreateDirectory(dirPath)
}
func OpenFile(filePath string, flag int, perm os.FileMode) (file *os.File, err error) {
	return FileHelper.OpenFile(filePath, flag, perm)
}

func DefOpenFile(filePath string) (file *os.File, err error) {
	return FileHelper.DefOpenFile(filePath)
}
func PathExist(filePath string) (exists bool) {
	return FileHelper.PathExist(filePath)
}
