package signature

import (
	"fmt"
	"github.com/soerenschneider/ssh-key-signer/internal"
	"github.com/soerenschneider/ssh-key-signer/pkg/ssh"

	log "github.com/rs/zerolog/log"
)

type Signer interface {
	SignPublicKey(publicKeyData string) (string, error)
	ReadCaCert() (string, error)
}

// KeyPod is a simple wrapper around a key (which is just a byte stream itself). This way, we decouple
// the implementation (file-based, memory, network, ..) and make it easily swap- and testable.
type KeyPod interface {
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

func (i *Issuer) SignHostCert(pubKey, signedKey KeyPod) error {
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

	newSignedKeyData, err := i.signerImpl.SignPublicKey(string(pubKeyData))
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
