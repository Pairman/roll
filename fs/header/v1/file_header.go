package v1

import (
	"bytes"
	"fmt"
	"os"

	"pnxlr.eu.org/roll/fs/header"
	"pnxlr.eu.org/roll/fs/util"
)

type FileHeader struct {
	PNGSect      PNGSect                      // PNG faking
	Size         int32                        // Header length
	Ver          header.FileHeaderVersionType // Format version
	FileSectSize int32                        // File information
	FileSect     FileSect                     //
	HashSectSize int32                        // File hashing
	HashSect     HashSect                     //
	CompSectSize int32                        // File compression
	CompSect     CompressionSect              //
	EncSectSize  int32                        // File encryption
	EncSect      EncryptionSect               //
}

func NewFileHeader(fileSect *FileSect, hashSect *HashSect,
	compressionSect *CompressionSect,
	encryptionSect *EncryptionSect) *FileHeader {
	pngSect := NewPNGSect()
	fileHeader := &FileHeader{
		PNGSect: *pngSect, Ver: header.FileHeaderVersion1,
		FileSectSize: int32(fileSect.Len()), FileSect: *fileSect,
		HashSectSize: int32(hashSect.Len()), HashSect: *hashSect,
		CompSectSize: int32(compressionSect.Len()), CompSect: *compressionSect,
		EncSectSize: int32(encryptionSect.Len()), EncSect: *encryptionSect,
	}
	fileHeader.Size = int32(fileHeader.Len())
	return fileHeader
}

func (s *FileHeader) String() string {
	return fmt.Sprintf("FileHeader{PNGSect: %v, Size: %v, Ver: %v, "+
		"FileSectSize: %v, FileSect: %v, HashSectSize: %v, HashSect: %v, "+
		"CompSectSize: %v, CompSect: %v, EncSectSize: %v, EncSect: %v}",
		s.PNGSect, s.Size, s.Ver,
		s.FileSectSize, s.FileSect, s.HashSectSize, s.HashSect,
		s.CompSectSize, s.CompSect, s.EncSectSize, s.EncSect)
}

func NewFileHeaderFromFile(f *os.File) *FileHeader {
	fs, _ := NewFileSect(f)
	hs, _ := NewHashSect(f, header.HashAlgoXXH3)
	cs, _ := NewCompressionSect(header.CompressionAlgoNone)
	es, _ := NewEncryptionSect(header.EncryptionAlgoNone)
	return NewFileHeader(fs, hs, cs, es)
}

func (fh *FileHeader) ToBytes() []byte {
	buf := bytes.Buffer{}
	buf.Write(fh.PNGSect.ToBytes())
	buf.Write(util.LiteralToBytes(fh.Size))
	buf.Write(util.LiteralToBytes(fh.Ver))
	buf.Write(util.LiteralToBytes(fh.FileSectSize))
	buf.Write(fh.FileSect.ToBytes())
	buf.Write(util.LiteralToBytes(fh.HashSectSize))
	buf.Write(fh.HashSect.ToBytes())
	buf.Write(util.LiteralToBytes(fh.CompSectSize))
	buf.Write(fh.CompSect.ToBytes())
	buf.Write(util.LiteralToBytes(fh.EncSectSize))
	buf.Write(fh.EncSect.ToBytes())
	return buf.Bytes()
}

func (fh *FileHeader) FromBytes(p []byte) error {
	fh.PNGSect.FromBytes(p[:SizePNGSectData])
	offs := SizePNGSectData
	fh.Size = util.LiteralFromBytes[int32](p[offs : offs+header.SizeInt32])
	offs += header.SizeInt32
	fh.Ver = util.LiteralFromBytes[header.FileHeaderVersionType](
		p[offs : offs+header.SizeFileHeaderVersionType])
	offs += header.SizeFileHeaderVersionType
	fh.FileSectSize = util.LiteralFromBytes[int32](
		p[offs : offs+header.SizeInt32])
	offs += header.SizeInt32
	fh.FileSect.FromBytes(p[offs : offs+int(fh.FileSectSize)])
	offs += int(fh.FileSectSize)
	fh.HashSectSize = util.LiteralFromBytes[int32](
		p[offs : offs+header.SizeInt32])
	offs += header.SizeInt32
	fh.HashSect.FromBytes(p[offs : offs+int(fh.HashSectSize)])
	offs += int(fh.HashSectSize)
	fh.CompSectSize = util.LiteralFromBytes[int32](
		p[offs : offs+header.SizeInt32])
	offs += header.SizeInt32
	fh.CompSect.FromBytes(p[offs : offs+int(fh.CompSectSize)])
	offs += int(fh.CompSectSize)
	fh.EncSectSize = util.LiteralFromBytes[int32](
		p[offs : offs+header.SizeInt32])
	offs += header.SizeInt32
	fh.EncSect.FromBytes(p[offs : offs+int(fh.EncSectSize)])
	offs += int(fh.EncSectSize)
	if offs != int(fh.Size) {
		return fmt.Errorf("file header length mismatch: %d, %d", fh.Size, offs)
	}
	return nil
}

func (fh *FileHeader) Len() int {
	return fh.PNGSect.Len() +
		header.SizeInt32 + header.SizeFileHeaderVersionType +
		header.SizeInt32 + fh.FileSect.Len() +
		header.SizeInt32 + fh.HashSect.Len() +
		header.SizeInt32 + fh.CompSect.Len() +
		header.SizeInt32 + fh.EncSect.Len()
}
