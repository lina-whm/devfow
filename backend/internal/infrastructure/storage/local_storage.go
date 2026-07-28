// Package storage provides file storage abstractions and local implementation.
package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Storage interface {
	Save(path string, reader io.Reader) error
	Get(path string) (io.ReadCloser, error)
	Delete(path string) error
}

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0750); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) resolvePath(path string) (string, error) {
	cleanPath := filepath.Clean(path)
	if cleanPath == "." || cleanPath == ".." {
		return "", fmt.Errorf("invalid path: %s", path)
	}
	fullPath := filepath.Join(s.basePath, cleanPath)
	if !strings.HasPrefix(fullPath, filepath.Clean(s.basePath)+string(filepath.Separator)) && fullPath != filepath.Clean(s.basePath) {
		return "", fmt.Errorf("path traversal detected: %s", path)
	}
	return fullPath, nil
}

func (s *LocalStorage) Save(path string, reader io.Reader) error {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() { _ = f.Close() }()
	_, err = io.Copy(f, reader)
	return err
}

func (s *LocalStorage) Get(path string) (io.ReadCloser, error) {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return f, nil
}

func (s *LocalStorage) Delete(path string) error {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return err
	}
	return os.Remove(fullPath)
}