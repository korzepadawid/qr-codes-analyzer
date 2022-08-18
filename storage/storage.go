package storage

import "context"

const (
	ImageMimeType = "image/png"
	ImageExt      = ".png"
)

type PutFileParams struct {
	Object      *[]byte
	StorageKey  string
	ContentType string
}

// FileStorage is responsible for performing bucket operations,
//such as create, read, update, delete
type FileStorage interface {

	// PutFile puts file into storage and returns nil, otherwise returns an error
	PutFile(ctx context.Context, params PutFileParams) error

	// DeleteFile deletes file from storage
	DeleteFile(ctx context.Context, storageKey string) error
}
