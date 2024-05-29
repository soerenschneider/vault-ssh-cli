package main

import (
	"errors"

	"github.com/hashicorp/vault/api"
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
	vaultAuth   vault.AuthMethod
	signingImpl *vault.SignatureClient
	issuer      *signature.Issuer
}

func dieOnErr(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}

func buildApp(config *config.Config) *app {
	app := &app{}
	var err error

	app.vaultClient, err = api.NewClient(vault.DeriveVaultConfig(config))
	dieOnErr(err, "could not build vault client")

	app.vaultAuth, err = buildAuthImpl(app.vaultClient, config)
	dieOnErr(err, "could not build auth strategy")

	vaultOpts := []vault.VaultOpts{
		vault.SshMountPath(config.VaultMountSsh),
		vault.VaultRole(config.VaultSshRole),
	}
	app.signingImpl, err = vault.NewVaultSigner(app.vaultClient, app.vaultAuth, vaultOpts...)
	dieOnErr(err, "could not build rotation client")

	renewStrategy, err := buildRenewalStrategy(config)
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
	return &signature.FileSink{FilePath: expanded}, nil
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

func buildAuthImpl(client *api.Client, conf *config.Config) (vault.AuthMethod, error) {
	token := conf.VaultToken
	if len(token) > 0 {
		return auth.NewTokenAuth(token)
	}

	if len(conf.VaultRoleId) > 0 && (len(conf.VaultSecretId) > 0 || len(conf.VaultSecretIdFile) > 0) {
		approleData := make(map[string]string)
		approleData[auth.KeyRoleId] = conf.VaultRoleId
		approleData[auth.KeySecretId] = conf.VaultSecretId
		approleData[auth.KeySecretIdFile] = conf.VaultSecretIdFile

		return auth.NewAppRoleAuth(client, approleData)
	}

	return auth.NewTokenImplicitAuth(), nil
}
