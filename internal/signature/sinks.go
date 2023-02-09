package signature

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

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
	return os.WriteFile(fs.FilePath, []byte(signedData), 0640)
}

func (fs *FileSink) CanWrite() error {
	dir := filepath.Dir(fs.FilePath)
	return unix.Access(dir, unix.W_OK)
}
