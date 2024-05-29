package signature

type SignUserKeyRequest struct {
	PublicKey  string `validation:"required"`
	Ttl        int    `validation:"gt=600"`
	Principals []string
	Extensions map[string]string
}

type SignHostKeyRequest struct {
	PublicKey  string `validation:"required"`
	Ttl        int    `validation:"gte=86400"`
	Principals []string
	Extensions map[string]string
}
