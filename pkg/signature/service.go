package signature

import (
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v3"
)

type CertType int

const (
	User CertType = 1 << iota
	Host CertType = 1 << iota
)

type Signer interface {
	SignUserKey(req SignatureRequest) (string, error)
	SignHostKey(req SignatureRequest) (string, error)
	ReadCaCert() (string, error)
}

// KeyStorage is a simple wrapper around a key (which is just a byte stream itself). This way, we decouple
// the implementation (file-based, memory, network, ..) and make it easily swap- and testable.
type KeyStorage interface {
	Read() ([]byte, error)
	Write(string) error
}

type SignatureService struct {
	signerImpl    Signer
	issueStrategy IssueStrategy
}

func NewSignatureService(signer Signer, issueStrategy IssueStrategy) (*SignatureService, error) {
	if signer == nil {
		return nil, errors.New("empty signer provided")
	}

	if issueStrategy == nil {
		return nil, errors.New("empty issue strategy implementation provided")
	}

	return &SignatureService{signerImpl: signer, issueStrategy: issueStrategy}, nil
}

func (i *SignatureService) SignUserCert(signRequest SignatureRequest, pubKey, signedKey KeyStorage) (*IssueResult, error) {
	pubKeyData, err := pubKey.Read()
	if err != nil {
		return nil, fmt.Errorf("could not read public key data: %w", err)
	}

	req := SignatureRequest{
		PublicKey:  string(pubKeyData),
		Ttl:        signRequest.Ttl,
		Principals: signRequest.Principals,
		Extensions: signRequest.Extensions,
		VaultRole:  signRequest.VaultRole,
	}

	signature := func() (string, error) {
		return i.signerImpl.SignUserKey(req)
	}

	return i.signCert(signedKey, signature, 3)
}

func (i *SignatureService) SignHostCert(req SignatureRequest, pubKey, signedKey KeyStorage) (*IssueResult, error) {
	pubKeyData, err := pubKey.Read()
	if err != nil {
		return nil, fmt.Errorf("could not read public key data: %w", err)
	}

	req.PublicKey = string(pubKeyData)

	signature := func() (string, error) {
		return i.signerImpl.SignHostKey(req)
	}

	return i.signCert(signedKey, signature, 3)
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

		if !i.issueStrategy.NeedsIssuing(&certInfo) {
			ret.Status = Noop
			return ret, nil
		}
	}

	var newSignedKeyData string
	op := func() error {
		newSignedKeyData, err = performSignature()
		return err
	}

	var backoffImpl backoff.BackOff = backoff.NewExponentialBackOff()
	backoffImpl = backoff.WithMaxRetries(backoffImpl, uint64(retries)) //nolint G115
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
