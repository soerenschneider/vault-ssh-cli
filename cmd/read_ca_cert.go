package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/ssh-key-signer/internal"
	"github.com/soerenschneider/ssh-key-signer/internal/signature"
	"github.com/soerenschneider/ssh-key-signer/internal/signature/vault"
	"strings"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func readCaCertCmd() *cobra.Command {
	var readCaCertCmd = &cobra.Command{
		Use:   "read-ca",
		Short: "Read ca cert from vault",
		Run:   readCaCertEntrypoint,
	}

	readCaCertCmd.PersistentFlags().StringP(FLAG_CA_FILE, "o", "", "Write the ca certificate to this output file")

	return readCaCertCmd
}

func readCaCertEntrypoint(ccmd *cobra.Command, args []string) {
	log.Info().Msgf("Starting up version %s (%s)", internal.BuildVersion, internal.CommitHash)
	config, err := config()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config")
	}
	config.Print()

	err = readCaCert(config)
	if err != nil {
		log.Fatal().Err(err).Msg("read-ca-crt unsuccessful")
	}
}

func readCaCert(config *Config) error {
	errors := config.ValidateCommon()
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		return fmt.Errorf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", "))
	}

	vaultClient, err := api.NewClient(getVaultConfig(config))
	if err != nil {
		return fmt.Errorf("could not build vault client: %v", err)
	}

	authImpl := vault.NewNoAuth()
	signingImpl, err := vault.NewVaultSigner(vaultClient, authImpl, config.VaultMountSsh, config.VaultSshRole)
	if err != nil {
		return fmt.Errorf("could not build vault impl: %v", err)
	}

	caCert, err := signingImpl.ReadCaCert()
	if err != nil {
		return err
	}

	var pod signature.KeyPod = &signature.BufferPod{Print: true}
	if len(config.CaFile) > 0 {
		pod = &signature.FsPod{FilePath: config.CaFile}
	}
	return pod.Write(caCert)
}
