package v1

import (
	"bytes"
	"fmt"
	"os"

	"pnxlr.eu.org/roll/fs/header"
	"pnxlr.eu.org/roll/fs/util"
)

type FileSect struct { // File information
	FileSize int64           // Length
	FileType header.FileType // File type
	Time     int64           // Modified time in milliseconds
	NameSize int16           // Name length
	Name     []byte          // Name
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
		FileType: header.FileTypeRegular,
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
	s.FileSize = util.LiteralFromBytes[int64](p[:header.SizeInt64])
	offs := header.SizeInt64
	s.FileType = util.LiteralFromBytes[header.FileType](
		p[offs : offs+header.SizeFileType])
	offs += header.SizeFileType
	s.Time = util.LiteralFromBytes[int64](p[offs : offs+header.SizeInt64])
	offs += header.SizeInt64
	s.NameSize = util.LiteralFromBytes[int16](p[offs : offs+header.SizeInt16])
	offs += header.SizeInt16
	if sectLen := s.Len(); len(p) != sectLen {
		return fmt.Errorf("buffer length mismatch: %d, %d", sectLen, len(p))
	}
	buff := bytes.Buffer{}
	buff.Write(p[offs:s.Len()])
	s.Name = buff.Bytes()
	return nil
}

func (s *FileSect) Len() int {
	return header.SizeInt64 + header.SizeFileType + header.SizeInt64 +
		header.SizeInt16 + int(s.NameSize)
}
