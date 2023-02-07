package ssh

import (
	"fmt"
	"math"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type CertInfo struct {
	Type        string
	Serial      uint64
	ValidAfter  time.Time
	ValidBefore time.Time
}

func (l *CertInfo) GetPercentage() float32 {
	total := l.ValidBefore.Sub(l.ValidAfter).Seconds()
	if total == 0 {
		return 0.
	}

	left := time.Until(l.ValidBefore).Seconds()
	return float32(math.Max(0, left*100/total))
}

func ParseCertData(pubKeyBytes []byte) (CertInfo, error) {
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	if err != nil {
		return CertInfo{}, fmt.Errorf("can not parse: %s", err)
	}

	cert, ok := pubKey.(*ssh.Certificate)
	if !ok {
		return CertInfo{}, fmt.Errorf("pub key is not a valid certificate: %s", err)
	}

	return CertInfo{
		Type:        cert.Type(),
		Serial:      cert.Serial,
		ValidBefore: time.Unix(int64(cert.ValidBefore), 0).UTC(),
		ValidAfter:  time.Unix(int64(cert.ValidAfter), 0).UTC(),
	}, nil
}

func ReadCertFromDisk(publicKeyFile string) (CertInfo, error) {
	bytes, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return CertInfo{}, fmt.Errorf("reading cert failed, can not read file: %v", err)
	}

	lifetime, err := ParseCertData(bytes)
	if err != nil {
		return CertInfo{}, fmt.Errorf("could not determine lifetime of cert: %v", err)
	}

	return lifetime, nil
}
