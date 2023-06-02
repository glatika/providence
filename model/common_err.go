package model

import "errors"

var (
	ErrRecordNotFound         = errors.New("record not found")
	ErrStockNotFound          = errors.New("stock not found")
	ErrPrivateKeyUnreadable   = errors.New("error reading private key")
	ErrFailedReadGRPCMetadata = errors.New("unable to read metadata")
	ErrFailedPrecondition     = errors.New("failed precondition parameter")
)
