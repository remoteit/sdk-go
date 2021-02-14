package tests

import (
	"testing"

	api "github.com/remoteit/sdk-go"
	apiContracts "github.com/remoteit/sdk-go/contracts"
)

func Test_Client_LoginUserPass(t *testing.T) {
	client := api.NewClient(apiContracts.DEFAULT_API_URL, APIKEY, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)

	authentication, errx := client.LoginWithPassword(USER, PASS)
	if errx != nil {
		t.Error(errx)
		t.FailNow()
	}

	if authentication.AuthHash == "" {
		t.Error("Missing AuthHash")
		t.FailNow()
	}
	if authentication.User == "" {
		t.Error("Missing User")
		t.FailNow()
	}
	if authentication.UserID == "" {
		t.Error("Missing UserID")
		t.FailNow()
	}
	if authentication.Token == "" {
		t.Error("Missing Token")
		t.FailNow()
	}
}
