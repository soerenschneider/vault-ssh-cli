package signature

import (
	"io/ioutil"
)

type BufferPod struct {
	Data []byte
}

func (b *BufferPod) Read() ([]byte, error) {
	return b.Data, nil
}

func (b *BufferPod) Write(signedData string) error {
	b.Data = []byte(signedData)
	return nil
}

type FsPod struct {
	FilePath string
}

func (fs *FsPod) Read() ([]byte, error) {
	return ioutil.ReadFile(fs.FilePath)
}

func (fs *FsPod) Write(signedData string) error {
	return ioutil.WriteFile(fs.FilePath, []byte(signedData), 0640)
}
