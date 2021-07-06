package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/qurname2/vaultAppRole"
	"net/http"
	"os"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	vaultAddr     = os.Getenv("VAULT_ADDR")
	vaultRoleID   = os.Getenv("VAULT_ROLE_ID")
	vaultSecretID = os.Getenv("VAULT_SECRET_ID")
)

func main() {
	config := &api.Config{
		Address:    vaultAddr,
		MaxRetries: 2,
		HttpClient: httpClient,
	}
	authParams := vaultAuth.AppRoleLogin{RoleID: vaultRoleID, SecretID: vaultSecretID}
	token, err := vaultAuth.AppRoleAuth(config, &authParams)
	if err != nil {
		fmt.Printf("error occured: %s", err)
	}
	fmt.Printf("vault-test, token: %s", token)
}
