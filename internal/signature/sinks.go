package signature

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"golang.org/x/sys/unix"
)

type AferoSink struct {
	fs       afero.Fs
	filePath string
}

func NewAferoSink(filePath string) (*AferoSink, error) {
	if len(filePath) == 0 {
		return nil, errors.New("empty filePath provided")
	}

	return &AferoSink{
		fs:       afero.NewOsFs(),
		filePath: filePath,
	}, nil
}

func (s *AferoSink) Read() ([]byte, error) {
	return afero.ReadFile(s.fs, s.filePath)
}

func (s *AferoSink) CanRead() error {
	_, err := afero.Exists(s.fs, s.filePath)
	return err
}

func (s *AferoSink) Write(signedData string) error {
	return afero.WriteFile(s.fs, s.filePath, []byte(signedData), 0640) // #nosec: G306
}

func (s *AferoSink) CanWrite() error {
	file, err := s.fs.OpenFile(s.filePath, os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

type BufferSink struct {
	Data  []byte
	Print bool
}

func (b *BufferSink) Read() ([]byte, error) {
	if len(b.Data) > 0 {
		return b.Data, nil
	}
	return nil, fmt.Errorf("empty buffer")
}

func (b *BufferSink) CanRead() error {
	if len(b.Data) > 0 {
		return nil
	}
	return fmt.Errorf("empty buffer")
}

func (b *BufferSink) Write(signedData string) error {
	b.Data = []byte(signedData)
	if b.Print {
		fmt.Println(signedData)
	}
	return nil
}

func (b *BufferSink) CanWrite() error {
	return nil
}

type FileSink struct {
	FilePath string
}

func (fs *FileSink) Read() ([]byte, error) {
	return os.ReadFile(fs.FilePath)
}

func (fs *FileSink) CanRead() error {
	_, err := os.Stat(fs.FilePath)
	return err
}

func (fs *FileSink) Write(signedData string) error {
	return os.WriteFile(fs.FilePath, []byte(signedData), 0640) // #nosec: G306
}

func (fs *FileSink) CanWrite() error {
	dir := filepath.Dir(fs.FilePath)
	return unix.Access(dir, unix.W_OK)
}
