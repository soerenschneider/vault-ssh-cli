package signature

import (
	"fmt"
	"testing"

	"github.com/soerenschneider/vault-ssh-cli/pkg/ssh"
)

const (
	signedDataContent  = "signed - so much secure"
	randomSshPublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBJxB104ZayRKBAcNRq+CqKcAI+bcbgJ4nIO7v0jWOZt"
	signedData         = `ssh-rsa-cert-v01@openssh.com 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`
)

type HappySignerDummy struct{}

func (s *HappySignerDummy) SignHostKey(publicKeyData string) (string, error) {
	return signedDataContent, nil
}

func (s *HappySignerDummy) SignUserKey(publicKeyData string) (string, error) {
	return signedDataContent, nil
}

func (s *HappySignerDummy) ReadCaCert() (string, error) {
	return "", nil
}

type SadSignerDummy struct{}

func (s *SadSignerDummy) SignHostKey(publicKeyData string) (string, error) {
	return "", fmt.Errorf("sad sad sad")
}

func (s *SadSignerDummy) SignUserKey(publicKeyData string) (string, error) {
	return "", fmt.Errorf("sad sad sad")
}

func (s *SadSignerDummy) ReadCaCert() (string, error) {
	return "", fmt.Errorf("sad sad sad")
}

func TestIssuer_SignHostCert(t *testing.T) {
	type fields struct {
		signerImpl  Signer
		refreshImpl ssh.RefreshSignatureStrategy
	}
	type args struct {
		pubKey    Sink
		signedKey Sink
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantErr           bool
		wantSignatureData string
	}{
		{
			name: "Happy path - no existing signed key",
			fields: fields{
				signerImpl:  &HappySignerDummy{},
				refreshImpl: ssh.NewSimpleStrategy(true),
			},
			args: args{
				pubKey:    &BufferSink{Data: []byte(randomSshPublicKey)},
				signedKey: &BufferSink{},
			},
			wantErr:           false,
			wantSignatureData: signedDataContent,
		},
		{
			name: "Happy path - existing signed key",
			fields: fields{
				signerImpl:  &HappySignerDummy{},
				refreshImpl: ssh.NewSimpleStrategy(true),
			},
			args: args{
				pubKey:    &BufferSink{Data: []byte(randomSshPublicKey)},
				signedKey: &BufferSink{Data: []byte(signedData)},
			},
			wantErr:           false,
			wantSignatureData: signedDataContent,
		},
		{
			name: "Happy path - existing signed key, renew strategy prohibits new signature",
			fields: fields{
				signerImpl:  &HappySignerDummy{},
				refreshImpl: ssh.NewSimpleStrategy(false),
			},
			args: args{
				pubKey:    &BufferSink{Data: []byte(randomSshPublicKey)},
				signedKey: &BufferSink{Data: []byte(signedData)},
			},
			wantErr:           false,
			wantSignatureData: signedData,
		},
		{
			name: "Error: sad signer",
			fields: fields{
				signerImpl:  &SadSignerDummy{},
				refreshImpl: ssh.NewSimpleStrategy(true),
			},
			args: args{
				pubKey:    &BufferSink{Data: []byte(randomSshPublicKey)},
				signedKey: &BufferSink{},
			},
			wantErr:           true,
			wantSignatureData: "",
		},
		{
			name: "Error: already signed certificate is garbage data",
			fields: fields{
				signerImpl:  &HappySignerDummy{},
				refreshImpl: ssh.NewSimpleStrategy(true),
			},
			args: args{
				pubKey:    &BufferSink{Data: []byte(randomSshPublicKey)},
				signedKey: &BufferSink{Data: []byte("garbage data, no ssh cert")},
			},
			wantErr:           true,
			wantSignatureData: "garbage data, no ssh cert",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Issuer{
				signerImpl:  tt.fields.signerImpl,
				refreshImpl: tt.fields.refreshImpl,
			}
			if err := i.SignHostCert(tt.args.pubKey, tt.args.signedKey); (err != nil) != tt.wantErr {
				t.Errorf("Issuer.SignHostCert() error = %v, wantErr %v", err, tt.wantErr)
			}

			signedData, _ := tt.args.signedKey.Read()
			if string(signedData) != tt.wantSignatureData {
				t.Errorf("Expected %s, got %s", signedDataContent, signedData)
			}
		})
	}
}
