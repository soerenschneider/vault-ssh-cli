package signature

type SignUserKeyRequest struct {
	PublicKey  string `validation:"required"`
	Ttl        string `validation:"gt=600"`
	Principals []string
	Extensions map[string]string
}

type SignHostKeyRequest struct {
	PublicKey  string `validation:"required"`
	Ttl        string `validation:"gte=86400"`
	Principals []string
	Extensions map[string]string
}
