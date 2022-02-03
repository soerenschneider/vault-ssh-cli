build:
	CGO_ENABLED=0 go build -o ssh-key-signer ./cmd

tests:
	true
