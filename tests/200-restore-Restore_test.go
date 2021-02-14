package tests

import (
	"testing"

	api "github.com/remoteit/sdk-go"
	apiContracts "github.com/remoteit/sdk-go/contracts"
)

func Test_Restore_Restore(t *testing.T) {
	client := api.NewClient(apiContracts.DEFAULT_API_URL, APIKEY, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)
	authentication, _ := client.LoginWithPassword(USER, PASS)

	restoreClient := api.NewRestoreClient(apiContracts.DEFAULT_API_RESTORE_URL, authentication.Token, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)

	configAsRaw, errx := restoreClient.Restore(DEVICEID, MACHINEID)
	if errx != nil {
		t.Error(errx)
		t.FailNow()
	}

	if configAsRaw != nil {
		t.Error(errx)
		t.FailNow()
	}
}
