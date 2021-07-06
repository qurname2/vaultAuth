package vaultAuth

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"os"
)

type AppRoleLogin struct {
	RoleID    string `json:"role_id,omitempty"`
	SecretID  string `json:"secret_id,omitempty"`
	LoginPath string `json:"login_path,omitempty"`
}

func AppRoleAuth(authConfig *api.Config, authParams *AppRoleLogin) (string, error) {
	client, err := api.NewClient(authConfig)
	if err != nil {
		return "", fmt.Errorf("error occurred with getting client for vault: %s", err.Error())
	}
	if authParams.RoleID == "" {
		authParams.RoleID = getEnv("VAULT_ROLE_ID", "")
	}
	if authParams.SecretID == "" {
		authParams.SecretID = getEnv("VAULT_SECRET_ID", "")
	}
	if authParams.LoginPath == "" {
		authParams.LoginPath = getEnv("VAULT_APPROLE_LOGIN_PATH", "/v1/auth/approle/login")
	}
	//create the token request
	request := client.NewRequest("POST", authParams.LoginPath)

	if err := request.SetJSONBody(authParams); err != nil {
		return "", fmt.Errorf("error occurred with SetJSONBody for a request: %s", err.Error())
	}
	//make the request
	resp, errRawRequest := client.RawRequest(request)
	if errRawRequest != nil {
		return "", fmt.Errorf("error occurred with RawRequest for a request: %s", errRawRequest.Error())
	}
	defer resp.Body.Close()

	//parse response
	secret, errParseSecret := api.ParseSecret(resp.Body)
	if errParseSecret != nil {
		return "", fmt.Errorf("error occurred with errParseSecret for a request: %s", errParseSecret.Error())
	}
	return secret.Auth.ClientToken, nil
}

// getEnv checks to see if an environment variable exists otherwise uses the default
//	env			: the name of the environment variable you are checking for
//	value		: the default value to return if the value is not there
func getEnv(env, value string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return value
}
