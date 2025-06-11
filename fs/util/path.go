package util

import (
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		if usr, err := user.Current(); err == nil {
			path = filepath.Join(usr.HomeDir, path[1:])
		}
	}
	path, _ = filepath.Abs(path)
	return path
}

func renameName(name string) string {
	ext := filepath.Ext(name)
	return name[:len(name)-len(ext)] + time.Now().Format(
		"_2006-01-02-15-04-05") + ext
}

func renamePath(path string) string {
	return filepath.Join(filepath.Dir(path), renameName(filepath.Base(path)))
}

func ComposePath(path, name string) (string, error) {
	path = expandPath(path)
	if info, err := os.Stat(path); os.IsNotExist(err) {
		if info, err := os.Stat(filepath.Dir(
			path)); err == nil && info.IsDir() {
			return path, nil
		}
		return "", nil
	} else if err != nil {
		return "", nil
	} else if info.IsDir() {
		path = filepath.Join(path, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
		} else if err != nil {
			return "", err
		} else {
			path = renamePath(path)
		}
	} else {
		path = renamePath(path)
	}
	return path, nil
}

func CreateFile(path string, size int64) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, file.Truncate(size)
}

func MoveFile(src string, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	} else if lerr, ok := err.(*os.LinkError); !ok ||
		lerr.Err != syscall.EXDEV {
		return err
	}

	ifile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer ifile.Close()

	ofile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer ofile.Close()

	if _, err := io.Copy(ofile, ifile); err != nil {
		return err
	}
	return os.Remove(src)
}
