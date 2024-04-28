package flyers

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

var (
	filenameEmptyErr             = errors.New("filename cannot be empty")
	defaultPerm      os.FileMode = 0777
)

// File 是对文件和目录的抽象
type File struct {
	path  string
	name  string
	IsDir bool

	option FileOptions
	logger *logrus.Entry
}

type FileOptions struct {
	Perm *os.FileMode
}

func NewDir(path string, ops ...FileOptions) *File {
	return NewFile(path, "", ops...)
}

func NewFile(path, name string, ops ...FileOptions) *File {
	var op FileOptions
	if len(ops) > 0 {
		op = ops[0]
	}
	if op.Perm == nil {
		op.Perm = &defaultPerm
	}
	return &File{
		path:   path,
		name:   name,
		option: op,
		logger: logrus.WithFields(logrus.Fields{
			"path": path,
			"name": name,
		}),
	}
}

func (f *File) GetFilePath() string {
	return path.Join(f.path, f.name)
}

func (f *File) PathExists() bool {
	if len(f.path) == 0 {
		return false
	}
	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		return false
	}
	return true
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

// ScanLine 逐行扫描
func (f *File) ScanLine(fn func(line string) bool) error {
	if !f.Exists() {
		return f.fileNotExistsErr()
	}
	file, err := f.openFile(os.O_RDONLY, 0666)
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

func (f *File) WriteLine(lines []string, isAtomic ...bool) error {
	if !f.PathExists() {
		err := os.MkdirAll(f.path, *f.option.Perm)
		if err != nil {
			return err
		}
	}
	file, err := f.openFile(os.O_WRONLY|os.O_CREATE, *f.option.Perm)
	defer file.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	errCnt := 0
	for _, line := range lines {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			if len(isAtomic) > 0 && isAtomic[0] {
				return err
			} else {
				errCnt += 1
			}
		}
	}
	if errCnt > 0 {
		f.logger.Warningf("write lines, error line nunber: %d", errCnt)
	}
	return writer.Flush()
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

func (f *File) WriteText(data string) error {
	fp := f.GetFilePath()
	if len(fp) == 0 {
		return filenameEmptyErr
	}
	return os.WriteFile(fp, []byte(data), *f.option.Perm)
}

func (f *File) openFile(flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(f.GetFilePath(), flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *File) fileNotExistsErr() error {
	return fmt.Errorf("file don't exists: %s", f.GetFilePath())
}
