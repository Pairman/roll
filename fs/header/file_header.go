package header

import (
	"bytes"
	"fmt"
	"os"

	"pnxlr.eu.org/roll/fs/util"
)

type FileHeaderVersionType int16

const (
	FileHeaderVersionUnknown FileHeaderVersionType = iota
	FileHeaderVersion1
)

type FileHeader struct {
	PNGSect      PNGSect               // PNG faking
	Size         int32                 // Header length
	Ver          FileHeaderVersionType // Format version
	FileSectSize int32                 // File information
	FileSect     FileSect              //
	HashSectSize int32                 // File hashing
	HashSect     HashSect              //
	CompSectSize int32                 // File compression
	CompSect     CompressionSect       //
	EncSectSize  int32                 // File encryption
	EncSect      EncryptionSect        //
}

func NewFileHeader(fileSect *FileSect, hashSect *HashSect,
	compressionSect *CompressionSect, encryptionSect *EncryptionSect) *FileHeader {
	pngSect := NewPNGSect()
	fileHeader := &FileHeader{
		PNGSect: *pngSect, Ver: FileHeaderVersion1,
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
	hs, _ := NewHashSect(f, 1)
	cs, _ := NewCompressionSect(CompressionAlgoNone)
	es, _ := NewEncryptionSect(EncryptionAlgoNone)
	return NewFileHeader(fs, hs, cs, es)
}

func (header *FileHeader) ToBytes() []byte {
	buf := bytes.Buffer{}
	buf.Write(header.PNGSect.ToBytes())
	buf.Write(util.LiteralToBytes(header.Size))
	buf.Write(util.LiteralToBytes(header.Ver))
	buf.Write(util.LiteralToBytes(header.FileSectSize))
	buf.Write(header.FileSect.ToBytes())
	buf.Write(util.LiteralToBytes(header.HashSectSize))
	buf.Write(header.HashSect.ToBytes())
	buf.Write(util.LiteralToBytes(header.CompSectSize))
	buf.Write(header.CompSect.ToBytes())
	buf.Write(util.LiteralToBytes(header.EncSectSize))
	buf.Write(header.EncSect.ToBytes())
	return buf.Bytes()
}

func (header *FileHeader) FromBytes(p []byte) error {
	header.PNGSect.FromBytes(p[:PNGSectDataSize])
	offs := PNGSectDataSize
	header.Size = util.LiteralFromBytes[int32](p[offs : offs+4])
	offs += 4
	header.Ver = util.LiteralFromBytes[FileHeaderVersionType](p[offs : offs+2])
	offs += 2
	header.FileSectSize = util.LiteralFromBytes[int32](p[offs : offs+4])
	offs += 4
	header.FileSect.FromBytes(p[offs : offs+int(header.FileSectSize)])
	offs += int(header.FileSectSize)
	header.HashSectSize = util.LiteralFromBytes[int32](p[offs : offs+4])
	offs += 4
	header.HashSect.FromBytes(p[offs : offs+int(header.HashSectSize)])
	offs += int(header.HashSectSize)
	header.CompSectSize = util.LiteralFromBytes[int32](p[offs : offs+4])
	offs += 4
	header.CompSect.FromBytes(p[offs : offs+int(header.CompSectSize)])
	offs += int(header.CompSectSize)
	header.EncSectSize = util.LiteralFromBytes[int32](p[offs : offs+4])
	offs += 4
	header.EncSect.FromBytes(p[offs : offs+int(header.EncSectSize)])
	offs += int(header.EncSectSize)
	if offs != int(header.Size) {
		return fmt.Errorf("file header length mismatch: %d, %d",
			header.Size, offs)
	}
	return nil
}

func (header *FileHeader) Len() int {
	return header.PNGSect.Len() + 4 + 2 +
		4 + header.FileSect.Len() +
		4 + header.HashSect.Len() +
		4 + header.CompSect.Len() +
		4 + header.EncSect.Len()
}
