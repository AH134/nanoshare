package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	path string
}

func NewLocalStorage(path string) (*LocalStorage, error) {
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, fmt.Errorf("creating storage dir: %w", err)
	}

	return &LocalStorage{
		path: path,
	}, nil
}

func (s *LocalStorage) Save(ctx context.Context, key string, r io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	dst := filepath.Join(s.path, key)
	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("local storage save: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return fmt.Errorf("local storage save: %w", err)
	}

	return nil
}

func (s *LocalStorage) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	dst := filepath.Join(s.path, key)
	file, err := os.Open(dst)
	if err != nil {
		return nil, fmt.Errorf("local storage open: %w", err)
	}

	return file, nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	dst := filepath.Join(s.path, key)
	if err := os.Remove(dst); err != nil {
		return fmt.Errorf("local storage delete: %w", err)
	}

	return nil
}
