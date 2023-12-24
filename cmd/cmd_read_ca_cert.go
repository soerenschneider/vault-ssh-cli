package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/vault-ssh-cli/internal"
	"github.com/soerenschneider/vault-ssh-cli/internal/config"
	"github.com/soerenschneider/vault-ssh-cli/internal/signature"
	"github.com/soerenschneider/vault-ssh-cli/internal/signature/vault"
	"github.com/soerenschneider/vault-ssh-cli/internal/signature/vault/auth"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func readCaCertCmd() *cobra.Command {
	var readCaCertCmd = &cobra.Command{
		Use:   "read-ca",
		Short: "Read ca cert from vault",
		PreRun: func(cmd *cobra.Command, args []string) {
			log.Info().Msgf("Starting up version %s (%s)", internal.BuildVersion, internal.CommitHash)
		},
		Run: readCaCertEntrypoint,
	}

	readCaCertCmd.PersistentFlags().StringP(config.FLAG_CA_FILE, "o", "", "Write the ca certificate to this output file")

	return readCaCertCmd
}

func readCaCertEntrypoint(ccmd *cobra.Command, args []string) {
	config, err := getConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config")
	}
	config.Print()

	if err := readCaCert(config); err != nil {
		log.Fatal().Err(err).Msg("read-ca-crt unsuccessful")
	}
}

func readCaCert(config *config.Config) error {
	errors := config.ValidateCommon()
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		return fmt.Errorf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", "))
	}

	vaultClient, err := api.NewClient(vault.DeriveVaultConfig(config))
	if err != nil {
		return fmt.Errorf("could not build vault client: %v", err)
	}

	vaultOpts := []vault.VaultOpts{vault.VaultRole(config.VaultSshRole)}
	signingImpl, err := vault.NewVaultSigner(vaultClient, &auth.NoAuth{}, vaultOpts...)
	if err != nil {
		return fmt.Errorf("could not build vault impl: %v", err)
	}

	caCert, err := signingImpl.ReadCaCert()
	if err != nil {
		return err
	}

	var pod signature.Sink = &signature.BufferSink{Print: true}
	if len(config.CaFile) > 0 {
		pod = &signature.FileSink{FilePath: config.CaFile}
	}
	return pod.Write(caCert)
}
