package header

import "strconv"

const (
	SizeInt8    = 1
	SizeInt16   = 2
	SizeInt32   = 4
	SizeInt64   = 8
	SizeUint8   = 1
	SizeUint16  = 2
	SizeUint32  = 4
	SizeUint64  = 8
	SizeFloat32 = 4
	SizeFloat64 = 8
)

type FileHeaderVersionType int16

const SizeFileHeaderVersionType = 2

const (
	FileHeaderVersionUnknown FileHeaderVersionType = iota
	FileHeaderVersion1
)

type CompressionAlgoType int8

const SizeCompressionAlgoType = 1

const (
	CompressionAlgoNone CompressionAlgoType = iota
	CompressionAlgoZSTD
)

func (a CompressionAlgoType) String() string {
	switch a {
	case CompressionAlgoNone:
		return "None"
	case CompressionAlgoZSTD:
		return "ZSTD"
	default:
		return strconv.Itoa(int(a))
	}
}

type EncryptionAlgoType int8

const SizeEncryptionAlgoType = 1

const (
	EncryptionAlgoNone EncryptionAlgoType = iota
	EncryptionAlgoAES256GCM
)

func (a EncryptionAlgoType) String() string {
	switch a {
	case EncryptionAlgoNone:
		return "None"
	case EncryptionAlgoAES256GCM:
		return "AES-256-GCM"
	default:
		return strconv.Itoa(int(a))
	}
}

type FileType int8

const SizeFileType = 1

const (
	FileTypeDirectory FileType = iota
	FileTypeRegular
	FileTypeSymlink
)

func (t FileType) String() string {
	switch t {
	case FileTypeDirectory:
		return "Directory"
	case FileTypeRegular:
		return "Regular"
	case FileTypeSymlink:
		return "Symlink"
	default:
		return strconv.Itoa(int(t))
	}
}

type HashAlgoType int8

const SizeHashAlgoType = 1

const (
	HashAlgoNone HashAlgoType = iota
	HashAlgoXXH3
)

func (a HashAlgoType) String() string {
	switch a {
	case HashAlgoNone:
		return "None"
	case HashAlgoXXH3:
		return "XXH3"
	default:
		return strconv.Itoa(int(a))
	}
}
