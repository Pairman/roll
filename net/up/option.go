package up

import headerV1 "pnxlr.eu.org/roll/fs/header/v1"

type UploadOptions struct {
	Verbose  bool
	Compress CompressionOptions
	Encrypt  EncryptionOptions
}

type CompressionOptions struct {
	On   bool
	Algo headerV1.CompressionAlgoType
}

type EncryptionOptions struct {
	On   bool
	Algo headerV1.EncryptionAlgoType
}
