package main

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-ssh-cli/internal/config"
	"github.com/soerenschneider/vault-ssh-cli/internal/signature"
	"github.com/soerenschneider/vault-ssh-cli/internal/signature/vault"
	"github.com/soerenschneider/vault-ssh-cli/internal/signature/vault/auth"
	"github.com/soerenschneider/vault-ssh-cli/pkg"
	"github.com/soerenschneider/vault-ssh-cli/pkg/ssh"
)

type app struct {
	vaultClient *api.Client
	vaultAuth   api.AuthMethod
	signingImpl *vault.SignatureClient
	issuer      *signature.Issuer
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

	vaultOpts := []vault.VaultOpts{
		vault.SshMountPath(conf.VaultMountSsh),
		vault.VaultRole(conf.VaultSshRole),
	}
	app.signingImpl, err = vault.NewVaultSigner(app.vaultClient.Logical(), vaultOpts...)
	dieOnErr(err, "could not build rotation client")

	renewStrategy, err := buildRenewalStrategy(conf)
	dieOnErr(err, "could not build signature refresh strategy")

	app.issuer, err = signature.NewIssuer(app.signingImpl, renewStrategy)
	dieOnErr(err, "could not build issuer")

	return app
}

type keys struct {
	pub  signature.Sink
	sign signature.Sink
}

func buildKeys(config *config.Config) *keys {
	var err error
	keys := &keys{}

	keys.pub, err = buildPublicKeySink(config)
	dieOnErr(err, "can't build sink to read public key from")

	err = keys.pub.CanRead()
	dieOnErr(err, "can not read from public key")

	keys.sign, err = buildSignedKeySink(config)
	dieOnErr(err, "can't build sink to write signature to")

	err = keys.sign.CanWrite()
	dieOnErr(err, "is not writable")

	return keys
}

func buildRenewalStrategy(config *config.Config) (ssh.RefreshSignatureStrategy, error) {
	if config.ForceNewSignature {
		return ssh.NewSimpleStrategy(true), nil
	}

	return ssh.NewPercentageStrategy(config.CertificateLifetimeThresholdPercentage)
}

func buildPublicKeySink(config *config.Config) (signature.Sink, error) {
	if nil == config {
		return nil, errors.New("empty config supplied")
	}

	if len(config.PublicKeyFile) == 0 {
		return nil, errors.New("no public key file supplied")
	}

	expanded := pkg.GetExpandedFile(config.PublicKeyFile)
	return signature.NewAferoSink(expanded)
}

func buildSignedKeySink(config *config.Config) (signature.Sink, error) {
	if nil == config {
		return nil, errors.New("empty config supplied")
	}

	if len(config.SignedKeyFile) == 0 {
		return nil, errors.New("no signed key file supplied")
	}

	expanded := pkg.GetExpandedFile(config.SignedKeyFile)
	return signature.NewAferoSink(expanded)
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
