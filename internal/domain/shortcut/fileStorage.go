package shortcut

import "io"

type FileStorage interface {
	Save(file io.ReadSeeker) (string, error)
	Delete(fileName string) error
}
