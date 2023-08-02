package signature

import (
	"fmt"

	"github.com/soerenschneider/vault-ssh-cli/internal"
	"github.com/soerenschneider/vault-ssh-cli/pkg/ssh"

	log "github.com/rs/zerolog/log"
)

type CertType int

const (
	User CertType = 1 << iota
	Host CertType = 1 << iota
)

type Signer interface {
	SignHostKey(publicKeyData string) (string, error)
	SignUserKey(publicKeyData string) (string, error)
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

func (i *Issuer) SignClientCert(pubKey, signedKey Sink) error {
	return i.signCert(pubKey, signedKey, User)
}

func (i *Issuer) SignHostCert(pubKey, signedKey Sink) error {
	return i.signCert(pubKey, signedKey, Host)
}

func (i *Issuer) signCert(pubKey, signedKey Sink, certType CertType) error {
	err := signedKey.CanWrite()
	if err != nil {
		return fmt.Errorf("not starting signing process, can't write to signedKeyPod: %v", err)
	}

	currentSignedKeyData, err := signedKey.Read()
	if err == nil {
		// on the first run the signed key data is not available, yet
		certInfo, err := ssh.ParseCertData(currentSignedKeyData)
		if err != nil {
			return fmt.Errorf("could not read certificate: %v", err)
		}
		log.Info().Msgf("Cert '%d' lifetime at %.1f%%, valid from %s to %s", certInfo.Serial, certInfo.GetPercentage(), certInfo.ValidAfter, certInfo.ValidBefore)

		updateCertMetrics(certInfo)
		if !i.refreshImpl.NeedsNewSignature(&certInfo) {
			log.Info().Msg("Re-signing certificate not necessary")
			return nil
		}
		log.Info().Msg("Requesting new signature for public key")
	}

	pubKeyData, err := pubKey.Read()
	if err != nil {
		return fmt.Errorf("could not read public key data: %v", err)
	}

	var newSignedKeyData string
	if certType == User {
		newSignedKeyData, err = i.signerImpl.SignUserKey(string(pubKeyData))
	} else {
		newSignedKeyData, err = i.signerImpl.SignHostKey(string(pubKeyData))
	}
	if err != nil {
		return fmt.Errorf("could not sign public key: %v", err)
	}
	log.Info().Msg("Received signed cert")

	certInfo, err := ssh.ParseCertData([]byte(newSignedKeyData))
	if err == nil {
		updateCertMetrics(certInfo)
	}

	log.Info().Msg("Writing signed cert")
	err = signedKey.Write(newSignedKeyData)
	if err != nil {
		return fmt.Errorf("can't write signed cert: %v", err)
	}
	return nil
}

func updateCertMetrics(certInfo ssh.CertInfo) {
	internal.MetricCertExpiry.Set(float64(certInfo.ValidBefore.Unix()))
	internal.MetricCertLifetimePercent.Set(float64(certInfo.GetPercentage()))
	internal.MetricCertLifetimeTotal.Set(certInfo.ValidBefore.Sub(certInfo.ValidAfter).Seconds())
}
