package tests

import (
	"testing"

	api "github.com/remoteit/sdk-go"
	apiContracts "github.com/remoteit/sdk-go/contracts"
)

func Test_Client_LoginUserAuthHash_Cache(t *testing.T) {
	client := api.NewClient(apiContracts.DEFAULT_API_URL, APIKEY, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)

	// get auth-hash + token
	authentication1, errx := client.LoginWithPassword(USER, PASS)
	if errx != nil {
		t.Error(errx)
		t.FailNow()
	}

	// get new token
	authentication2, errx := client.LoginWithAuthHash(USER, authentication1.AuthHash)

	// get cached token
	authentication3, errx := client.LoginWithAuthHash(USER, authentication2.AuthHash)

	if authentication3.AuthHash == "" {
		t.Error("Missing AuthHash")
		t.FailNow()
	}
	if authentication3.User == "" {
		t.Error("Missing User")
		t.FailNow()
	}
	if authentication3.UserID == "" {
		t.Error("Missing UserID")
		t.FailNow()
	}
	if authentication3.Token == "" || authentication3.Token != authentication2.Token {
		t.Error("Missing Token")
		t.FailNow()
	}
}
