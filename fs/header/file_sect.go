package header

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"pnxlr.eu.org/roll/fs/util"
)

type FileType int8

const (
	FileTypeDirectory FileType = iota
	FileTypeRegular
	FileTypeSymlink
)

func (a FileType) String() string {
	switch a {
	case FileTypeDirectory:
		return "Directory"
	case FileTypeRegular:
		return "Regular"
	case FileTypeSymlink:
		return "Symlink"
	default:
		return strconv.Itoa(int(a))
	}
}

type FileSect struct { // File information
	FileSize int64    // Length
	FileType FileType // File type
	Time     int64    // Modified time in milliseconds
	NameSize int16    // Name length
	Name     []byte   // Name
}

func (s FileSect) String() string {
	return fmt.Sprintf("FileSect{FileSize: %v, FileType: %v, Time: %v, "+
		"NameSize: %v, Name: %v}",
		s.FileSize, s.FileType, s.Time, s.NameSize, string(s.Name))
}

func NewFileSect(file *os.File) (*FileSect, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	name := []byte(fileInfo.Name())
	return &FileSect{
		FileSize: fileInfo.Size(),
		FileType: FileTypeRegular,
		Time:     fileInfo.ModTime().UnixMilli(),
		NameSize: int16(len(name)), Name: name,
	}, nil
}

func (s FileSect) ToBytes() []byte {
	buf := bytes.Buffer{}
	buf.Write(util.LiteralToBytes(s.FileSize))
	buf.Write(util.LiteralToBytes(s.FileType))
	buf.Write(util.LiteralToBytes(s.Time))
	buf.Write(util.LiteralToBytes(s.NameSize))
	buf.Write(s.Name)
	return buf.Bytes()
}

func (s *FileSect) FromBytes(p []byte) error {
	s.FileSize = util.LiteralFromBytes[int64](p[:8])
	s.FileType = util.LiteralFromBytes[FileType](p[8:9])
	s.Time = util.LiteralFromBytes[int64](p[9:17])
	s.NameSize = util.LiteralFromBytes[int16](p[17:19])
	if sectLen := s.Len(); len(p) != sectLen {
		return fmt.Errorf("buffer length mismatch: %d, %d", sectLen, len(p))
	}
	buff := bytes.Buffer{}
	buff.Write(p[19:s.Len()])
	s.Name = buff.Bytes()
	return nil
}

func (s *FileSect) Len() int {
	return 19 + int(s.NameSize)
}
