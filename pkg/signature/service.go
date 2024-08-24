package signature

import (
	"fmt"

	"github.com/cenkalti/backoff/v3"
	"github.com/soerenschneider/vault-ssh-cli/internal/config"
)

type CertType int

const (
	User CertType = 1 << iota
	Host CertType = 1 << iota
)

type Signer interface {
	SignUserKey(req SignUserKeyRequest) (string, error)
	SignHostKey(req SignHostKeyRequest) (string, error)
	ReadCaCert() (string, error)
}

// KeyStorage is a simple wrapper around a key (which is just a byte stream itself). This way, we decouple
// the implementation (file-based, memory, network, ..) and make it easily swap- and testable.
type KeyStorage interface {
	Read() ([]byte, error)
	Write(string) error
}

type SignatureService struct {
	signerImpl  Signer
	refreshImpl RefreshSignatureStrategy
}

func NewSignatureService(signer Signer, refresh RefreshSignatureStrategy) (*SignatureService, error) {
	return &SignatureService{signerImpl: signer, refreshImpl: refresh}, nil
}

func (i *SignatureService) SignUserCert(conf *config.Config, pubKey, signedKey KeyStorage) (*IssueResult, error) {
	pubKeyData, err := pubKey.Read()
	if err != nil {
		return nil, fmt.Errorf("could not read public key data: %w", err)
	}

	req := SignUserKeyRequest{
		PublicKey:  string(pubKeyData),
		Ttl:        conf.Ttl,
		Principals: conf.Principals,
		Extensions: conf.Extensions,
	}

	signature := func() (string, error) {
		return i.signerImpl.SignUserKey(req)
	}

	return i.signCert(signedKey, signature, conf.Retries)
}

func (i *SignatureService) SignHostCert(conf *config.Config, pubKey, signedKey KeyStorage) (*IssueResult, error) {
	pubKeyData, err := pubKey.Read()
	if err != nil {
		return nil, fmt.Errorf("could not read public key data: %w", err)
	}

	req := SignHostKeyRequest{
		PublicKey:  string(pubKeyData),
		Ttl:        conf.Ttl,
		Principals: conf.Principals,
		Extensions: conf.Extensions,
	}

	signature := func() (string, error) {
		return i.signerImpl.SignHostKey(req)
	}

	return i.signCert(signedKey, signature, conf.Retries)
}

func (i *SignatureService) signCert(signedKey KeyStorage, performSignature func() (string, error), retries int) (*IssueResult, error) {
	ret := &IssueResult{
		Status: Unknown,
	}

	currentSignedKeyData, err := signedKey.Read()
	if err == nil {
		// on the first run the signed key data is not available, yet
		certInfo, err := ParseCertData(currentSignedKeyData)
		ret.ExistingCert = &certInfo
		if err != nil {
			return nil, fmt.Errorf("could not read certificate: %v", err)
		}

		if !i.refreshImpl.NeedsNewSignature(&certInfo) {
			ret.Status = Noop
			return ret, nil
		}
	}

	var newSignedKeyData string
	op := func() error {
		newSignedKeyData, err = performSignature()
		return err
	}

	var backoffImpl backoff.BackOff
	backoffImpl = backoff.NewExponentialBackOff()
	backoffImpl = backoff.WithMaxRetries(backoffImpl, uint64(retries))
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return ret, fmt.Errorf("could not sign public key: %w", err)
	}

	certInfo, err := ParseCertData([]byte(newSignedKeyData))
	if err != nil {
		return ret, fmt.Errorf("could not parse received cert data: %w", err)
	}
	ret.IssuedCert = &certInfo
	ret.Status = Issued

	return ret, signedKey.Write(newSignedKeyData)
}
