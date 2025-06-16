package up

import "pnxlr.eu.org/roll/fs/header"

type UploadOptions struct {
	Verbose  bool
	Compress CompressionOptions
	Encrypt  EncryptionOptions
}

type CompressionOptions struct {
	On   bool
	Algo header.CompressionAlgoType
}

type EncryptionOptions struct {
	On   bool
	Algo header.EncryptionAlgoType
}
