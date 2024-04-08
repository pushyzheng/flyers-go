package flyers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
)

var (
	filenameEmptyErr = errors.New("filename cannot be empty")
)

// File 是对文件和目录的抽象
type File struct {
	path  string
	name  string
	IsDir bool
}

func NewFile(path, name string) *File {
	return &File{
		path: path,
		name: name,
	}
}

func NewDir(path string) *File {
	return &File{
		path:  path,
		IsDir: true,
	}
}

func (f *File) GetFilePath() string {
	return path.Join(f.path, f.name)
}

func (f *File) Exists() bool {
	if len(f.GetFilePath()) == 0 {
		return false
	}
	if _, err := os.Stat(f.GetFilePath()); os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *File) Remove() error {
	if !f.Exists() {
		return f.fileNotExistsErr()
	}
	return os.Remove(f.GetFilePath())
}

func (f *File) MkDirIfNotExists() error {
	return os.MkdirAll(f.path, 0755)
}

func (f *File) ScanLine(fn func(line string) bool) error {
	if !f.Exists() {
		return f.fileNotExistsErr()
	}
	file, err := os.Open(f.GetFilePath())
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		continueF := fn(line)
		if !continueF {
			break
		}
	}
	return nil
}

func (f *File) ReadText() (string, error) {
	fp := f.GetFilePath()
	if len(fp) == 0 {
		return "", filenameEmptyErr
	}
	if !f.Exists() {
		return "", f.fileNotExistsErr()
	}
	b, err := os.ReadFile(fp)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (f *File) fileNotExistsErr() error {
	return fmt.Errorf("file don't exists: %s", f.GetFilePath())
}

func (f *File) WriteText(data string) error {
	fp := f.GetFilePath()
	if len(fp) == 0 {
		return filenameEmptyErr
	}
	return os.WriteFile(fp, []byte(data), 0666)
}
