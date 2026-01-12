// Package storage provides a wrapper around omnistorage backends.
package storage

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/grokify/omnistorage"

	// Register backends
	_ "github.com/grokify/omnistorage-github/backend/github"
	_ "github.com/grokify/omnistorage/backend/file"
	_ "github.com/grokify/omnistorage/backend/memory"
)

// Storage wraps an omnistorage.Backend with ChatHub-specific operations.
type Storage struct {
	backend omnistorage.Backend
	folder  string
}

// New creates a new Storage instance.
func New(backend omnistorage.Backend, folder string) *Storage {
	return &Storage{
		backend: backend,
		folder:  folder,
	}
}

// NewFromConfig creates a Storage from backend name and config.
func NewFromConfig(backendName string, config map[string]string, folder string) (*Storage, error) {
	backend, err := omnistorage.Open(backendName, config)
	if err != nil {
		return nil, fmt.Errorf("failed to open backend %s: %w", backendName, err)
	}
	return New(backend, folder), nil
}

// Close closes the underlying backend.
func (s *Storage) Close() error {
	return s.backend.Close()
}

// Save writes content to a path.
func (s *Storage) Save(ctx context.Context, filePath string, content []byte) error {
	w, err := s.backend.NewWriter(ctx, filePath)
	if err != nil {
		return fmt.Errorf("failed to create writer for %s: %w", filePath, err)
	}

	_, err = w.Write(content)
	if err != nil {
		w.Close()
		return fmt.Errorf("failed to write to %s: %w", filePath, err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close writer for %s: %w", filePath, err)
	}

	return nil
}

// Read reads content from a path.
func (s *Storage) Read(ctx context.Context, filePath string) ([]byte, error) {
	r, err := s.backend.NewReader(ctx, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", filePath, err)
	}
	defer r.Close()

	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	return content, nil
}

// List lists files with a prefix.
func (s *Storage) List(ctx context.Context, prefix string) ([]string, error) {
	files, err := s.backend.List(ctx, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list %s: %w", prefix, err)
	}
	return files, nil
}

// ListConversations lists all conversation files in the storage folder.
func (s *Storage) ListConversations(ctx context.Context) ([]string, error) {
	return s.List(ctx, s.folder)
}

// ListBySource lists conversations from a specific source.
func (s *Storage) ListBySource(ctx context.Context, source string) ([]string, error) {
	prefix := path.Join(s.folder, source)
	files, err := s.List(ctx, prefix)
	if err != nil {
		return nil, err
	}

	// Filter to only .md files
	var mdFiles []string
	for _, f := range files {
		if strings.HasSuffix(f, ".md") {
			mdFiles = append(mdFiles, f)
		}
	}
	return mdFiles, nil
}

// Delete removes a file.
func (s *Storage) Delete(ctx context.Context, filePath string) error {
	if err := s.backend.Delete(ctx, filePath); err != nil {
		return fmt.Errorf("failed to delete %s: %w", filePath, err)
	}
	return nil
}

// Exists checks if a file exists.
func (s *Storage) Exists(ctx context.Context, filePath string) (bool, error) {
	return s.backend.Exists(ctx, filePath)
}

// Folder returns the configured root folder.
func (s *Storage) Folder() string {
	return s.folder
}

// Backend returns the underlying omnistorage.Backend.
func (s *Storage) Backend() omnistorage.Backend {
	return s.backend
}
