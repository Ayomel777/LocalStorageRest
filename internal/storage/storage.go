package storage

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrBadPath  = errors.New("bad path")
	ErrNotFound = errors.New("invalid path")
)

type Storage struct {
	base string `json:"base"`
}

func NewStorage(base string) (*Storage, error) {
	absPath, err := filepath.Abs(base)
	if err != nil {
		return nil, err
	}

	return &Storage{
		base: absPath,
	}, os.MkdirAll(absPath, 0755)
}

func (s *Storage) safePath(p string) (string, error) {
	abs, err := filepath.Abs(filepath.Join(s.base, filepath.FromSlash(p)))
	if err != nil {
		return "", ErrBadPath
	}
	if abs != s.base && !strings.HasPrefix(abs, s.base+string(os.PathSeparator)) {
		return "", ErrBadPath
	}
	return abs, nil
}

func (s *Storage) Read(path string) (io.ReadCloser, error) {
	abs, err := s.safePath(path)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(abs, os.O_RDONLY, 0755)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}

	return f, nil
}

func (s *Storage) Write(path string, r io.Reader) error {
	abs, err := s.safePath(path)
	if err != nil {
		return err
	}
	_ = os.MkdirAll(filepath.Dir(abs), 0755)
	//f, err := os.OpenFile(abs, os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	log.Printf("Create error: %v", err)
	//	return err
	//}
	//defer f.Close()
	f, err := os.Create(abs)
	if err != nil {
		log.Printf("Create error: %v", err)
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

func (s *Storage) Append(path string, r io.Reader) error {
	abs, err := s.safePath(path)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(abs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

func (s *Storage) Delete(path string) error {
	abs, err := s.safePath(path)
	if err != nil {
		return err
	}

	err = os.Remove(abs)
	if os.IsNotExist(err) {
		return ErrNotFound
	}
	return err
}

func (s *Storage) Copy(srcPath, destPath string) error {
	src, err := s.safePath(srcPath)
	if err != nil {
		return err
	}

	dst, err := s.safePath(destPath)
	if err != nil {
		return err
	}

	fileIn, err := os.OpenFile(src, os.O_RDONLY, 0755)
	if os.IsNotExist(err) {
		return ErrNotFound
	}
	defer fileIn.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	fileOut, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileOut, fileIn)
	return err
}

func (s *Storage) Move(src, dst string) error {
	absSrc, err := s.safePath(src)
	if err != nil {
		return err
	}
	absDst, err := s.safePath(dst)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absSrc); os.IsNotExist(err) {
		return ErrNotFound
	}
	if err = os.MkdirAll(filepath.Dir(absDst), 0755); err != nil {
		return err
	}
	return os.Rename(absSrc, absDst)
}
