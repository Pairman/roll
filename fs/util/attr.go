package util

import (
	"os"
	"time"
)

func SetFileMTime(file *os.File, ms int64) error {
	t := time.UnixMilli(ms)
	return os.Chtimes(file.Name(), t, t)
}
