package main

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-ssh-cli/internal"
	"github.com/soerenschneider/vault-ssh-cli/internal/config"
	"github.com/soerenschneider/vault-ssh-cli/internal/vault"
	"github.com/soerenschneider/vault-ssh-cli/internal/vault/auth"
	"github.com/soerenschneider/vault-ssh-cli/pkg"
	"github.com/soerenschneider/vault-ssh-cli/pkg/signature"
)

type app struct {
	vaultClient *api.Client
	vaultAuth   api.AuthMethod

	signatureClient  *signature.SignatureClient
	signatureService *signature.SignatureService
}

func dieOnErr(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}

func buildApp(conf *config.Config) *app {
	app := &app{}
	var err error

	app.vaultClient, err = api.NewClient(vault.FromConfig(conf))
	dieOnErr(err, "could not build vault client")

	app.vaultAuth, err = buildAuthImpl(app.vaultClient, conf)
	dieOnErr(err, "could not build auth strategy")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = app.vaultAuth.Login(ctx, app.vaultClient)
	dieOnErr(err, "could not login to vault")

	vaultOpts := []signature.VaultOpts{
		signature.WithSshMountPath(conf.VaultMountSsh),
		signature.WithVaultRole(conf.VaultSshRole),
	}
	app.signatureClient, err = signature.NewVaultSigner(app.vaultClient.Logical(), vaultOpts...)
	dieOnErr(err, "could not build rotation client")

	renewStrategy, err := buildRenewalStrategy(conf)
	dieOnErr(err, "could not build signature refresh strategy")

	app.signatureService, err = signature.NewSignatureService(app.signatureClient, renewStrategy)
	dieOnErr(err, "could not build issuer")

	return app
}

type keys struct {
	pub  signature.KeyStorage
	sign signature.KeyStorage
}

func buildKeys(config *config.Config) *keys {
	var err error
	keys := &keys{}

	keys.pub, err = buildPublicKeyStorage(config)
	dieOnErr(err, "can't build sink to read public key from")

	keys.sign, err = buildSignedKeyStorage(config)
	dieOnErr(err, "can't build sink to write signature to")

	return keys
}

func buildRenewalStrategy(config *config.Config) (signature.RefreshSignatureStrategy, error) {
	if config.ForceNewSignature {
		return signature.NewSimpleStrategy(true), nil
	}

	return signature.NewPercentageStrategy(config.CertificateLifetimeThresholdPercentage)
}

func buildPublicKeyStorage(config *config.Config) (signature.KeyStorage, error) {
	if nil == config {
		return nil, errors.New("empty config supplied")
	}

	if len(config.PublicKeyFile) == 0 {
		return nil, errors.New("no public key file supplied")
	}

	expanded := pkg.GetExpandedFile(config.PublicKeyFile)
	return internal.NewAferoSink(expanded)
}

func buildSignedKeyStorage(config *config.Config) (signature.KeyStorage, error) {
	if nil == config {
		return nil, errors.New("empty config supplied")
	}

	if len(config.SignedKeyFile) == 0 {
		return nil, errors.New("no signed key file supplied")
	}

	expanded := pkg.GetExpandedFile(config.SignedKeyFile)
	return internal.NewAferoSink(expanded)
}

func buildAuthImpl(client *api.Client, conf *config.Config) (api.AuthMethod, error) {
	if len(conf.VaultRoleId) > 0 && (len(conf.VaultSecretId) > 0 || len(conf.VaultSecretIdFile) > 0) {
		secretId := &approle.SecretID{}
		if len(conf.VaultSecretIdFile) > 0 {
			secretId.FromFile = conf.VaultSecretIdFile
		} else {
			secretId.FromString = conf.VaultSecretId
		}
		return approle.NewAppRoleAuth(conf.VaultRoleId, secretId)

	}

	return &auth.NoAuth{}, nil
}
