package signature

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

type BufferPod struct {
	Data []byte
}

func (b *BufferPod) Read() ([]byte, error) {
	if len(b.Data) > 0 {
		return b.Data, nil
	}
	return nil, fmt.Errorf("empty buffer")
}

func (b *BufferPod) CanRead() error {
	if len(b.Data) > 0 {
		return nil
	}
	return fmt.Errorf("empty buffer")
}

func (b *BufferPod) Write(signedData string) error {
	b.Data = []byte(signedData)
	return nil
}

func (b *BufferPod) CanWrite() error {
	return nil
}

type FsPod struct {
	FilePath string
}

func (fs *FsPod) Read() ([]byte, error) {
	return ioutil.ReadFile(fs.FilePath)
}

func (fs *FsPod) CanRead() error {
	_, err := os.Stat(fs.FilePath)
	return err
}

func (fs *FsPod) Write(signedData string) error {
	return ioutil.WriteFile(fs.FilePath, []byte(signedData), 0640)
}

func (fs *FsPod) CanWrite() error {
	dir := filepath.Dir(fs.FilePath)
	return unix.Access(dir, unix.W_OK)
}
