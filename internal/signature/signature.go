package signature

import (
	"fmt"
	"github.com/soerenschneider/ssh-key-signer/internal"
	"github.com/soerenschneider/ssh-key-signer/pkg/ssh"

	log "github.com/rs/zerolog/log"
)

type Signer interface {
	SignPublicKey(publicKeyData string) (string, error)
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
	if err == nil {
		// on the first run the signed key data is not available, yet
		certInfo, err := ssh.ParseCertData(pubKeyData)
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

	pubKeyData, err = pubKey.Read()
	if err != nil {
		return fmt.Errorf("could not read public key data: %v", err)
	}

	signedData, err := i.signerImpl.SignPublicKey(string(pubKeyData))
	if err != nil {
		return fmt.Errorf("could not sign public key: %v", err)
	}

	certInfo, err := ssh.ParseCertData([]byte(signedData))
	if err == nil {
		updateCertMetrics(certInfo)
	}

	log.Info().Msg("Writing signed cert")
	err = signedKey.Write(signedData)
	if err != nil {
		return fmt.Errorf("could not write signed public key to pod: %v", err)
	}

	return nil
}

func updateCertMetrics(certInfo ssh.CertInfo) {
	internal.MetricCertExpiry.Set(float64(certInfo.ValidBefore.Unix()))
	internal.MetricCertLifetime.Set(float64(certInfo.GetPercentage()))
}
