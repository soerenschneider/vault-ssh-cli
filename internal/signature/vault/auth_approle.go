package vault

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/vault/api"
)

type AppRoleAuth struct {
	client *api.Client

	approleMount string
	loginData    map[string]string
}

func NewAppRoleAuth(client *api.Client, loginData map[string]string) (*AppRoleAuth, error) {
	return &AppRoleAuth{
		client:    client,
		loginData: loginData,
		// TODO: supply variable
		approleMount: "approle",
	}, nil
}

func (t *AppRoleAuth) getLoginTuple() (string, string, error) {
	var roleId, secretId string
	val, ok := t.loginData["role_id"]
	if ok && len(val) > 0 {
		roleId = val
	}

	val, ok = t.loginData["secret_id"]
	if ok && len(val) > 0 {
		secretId = val
	}

	val, ok = t.loginData["secret_id_file"]
	if ok && len(val) > 0 {
		data, err := ioutil.ReadFile(val)
		if err != nil {
			return "", "", fmt.Errorf("could not read secret_id from file '%s': %v", val, err)
		}
		secretId = string(data)
	}

	return roleId, secretId, nil
}

func (t *AppRoleAuth) Authenticate() (string, error) {
	roleId, secretId, err := t.getLoginTuple()
	if err != nil {
		return "", fmt.Errorf("could not get login data: %v", err)
	}

	path := fmt.Sprintf("auth/%s/login", t.approleMount)
	data := map[string]interface{}{
		"role_id":   roleId,
		"secret_id": secretId,
	}
	secret, err := t.client.Logical().Write(path, data)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}
