package storage

import (
	"errors"
	"fmt"
	"intellisys-backend/internal/infra/config"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(cfg *config.Config) *LocalStorage {
	return &LocalStorage{
		BasePath: cfg.FalePath.Path,
	}
}

func (s *LocalStorage) Save(file io.ReadSeeker) (string, error) {
	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buffer[:n])

	ext := getExtension(contentType)
	if ext == "" {
		return "", errors.New("arquivo inválido")
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	newName := uuid.New().String() + ext

	err = os.MkdirAll(s.BasePath, os.ModePerm)
	if err != nil {
		return "", err
	}

	path := filepath.Join(s.BasePath, newName)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return newName, nil
}

func (s *LocalStorage) Delete(fileName string) error {
	if fileName == "" {
		return nil
	}

	fileName = filepath.Base(fileName)

	fullPath := filepath.Join(s.BasePath, fileName)

	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("erro ao deletar imagem: %w", err)
	}

	return nil
}

func getExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
