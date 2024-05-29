package signature

import (
	"fmt"
	"time"

	"github.com/soerenschneider/vault-ssh-cli/internal"
	"github.com/soerenschneider/vault-ssh-cli/internal/config"
	"github.com/soerenschneider/vault-ssh-cli/pkg/ssh"

	log "github.com/rs/zerolog/log"
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

// Sink is a simple wrapper around a key (which is just a byte stream itself). This way, we decouple
// the implementation (file-based, memory, network, ..) and make it easily swap- and testable.
type Sink interface {
	Read() ([]byte, error)
	CanRead() error
	Write(string) error
	CanWrite() error
}

type Issuer struct {
	signerImpl  Signer
	refreshImpl ssh.RefreshSignatureStrategy
}

func NewIssuer(signer Signer, refresh ssh.RefreshSignatureStrategy) (*Issuer, error) {
	return &Issuer{signerImpl: signer, refreshImpl: refresh}, nil
}

func (i *Issuer) SignUserCert(conf *config.Config, pubKey, signedKey Sink) error {
	pubKeyData, err := pubKey.Read()
	if err != nil {
		return fmt.Errorf("could not read public key data: %w", err)
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

	return i.signCert(signedKey, signature)
}

func (i *Issuer) SignHostCert(conf *config.Config, pubKey, signedKey Sink) error {
	pubKeyData, err := pubKey.Read()
	if err != nil {
		return fmt.Errorf("could not read public key data: %w", err)
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

	return i.signCert(signedKey, signature)

}

func (i *Issuer) signCert(signedKey Sink, performSignature func() (string, error)) error {
	if err := signedKey.CanWrite(); err != nil {
		return fmt.Errorf("not starting signing process, can't write to signedKeyPod: %v", err)
	}

	currentSignedKeyData, err := signedKey.Read()
	if err == nil {
		// on the first run the signed key data is not available, yet
		certInfo, err := ssh.ParseCertData(currentSignedKeyData)
		if err != nil {
			return fmt.Errorf("could not read certificate: %v", err)
		}
		log.Info().Msgf("Cert '%d' lifetime at %.1f%%, valid from '%s' to '%s'", certInfo.Serial, certInfo.GetPercentage(), certInfo.ValidAfter, certInfo.ValidBefore)

		updateCertMetrics(certInfo)
		if !i.refreshImpl.NeedsNewSignature(&certInfo) {
			log.Info().Msg("Signing certificate not necessary")
			return nil
		}
		log.Info().Msg("Requesting new signature for public key")
	}

	newSignedKeyData, err := performSignature()
	if err != nil {
		return fmt.Errorf("could not sign public key: %w", err)
	}

	certInfo, err := ssh.ParseCertData([]byte(newSignedKeyData))
	if err != nil {
		return fmt.Errorf("could not parse received cert data: %w", err)
	}
	log.Info().Msgf("Received signed SSH cert, valid until %s (%v)", certInfo.ValidBefore, time.Until(certInfo.ValidBefore))
	updateCertMetrics(certInfo)

	return signedKey.Write(newSignedKeyData)
}

func updateCertMetrics(certInfo ssh.CertInfo) {
	internal.MetricCertExpiry.Set(float64(certInfo.ValidBefore.Unix()))
	internal.MetricCertLifetimePercent.Set(float64(certInfo.GetPercentage()))
	internal.MetricCertLifetimeTotal.Set(certInfo.ValidBefore.Sub(certInfo.ValidAfter).Seconds())
}
