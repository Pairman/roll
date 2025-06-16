package header

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
