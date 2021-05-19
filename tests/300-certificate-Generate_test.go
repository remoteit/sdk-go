package tests

import (
	"testing"

	api "github.com/remoteit/sdk-go"
	apiContracts "github.com/remoteit/sdk-go/contracts"
)

func Test_Certificate_Generate(t *testing.T) {
	client := api.NewClient(apiContracts.DEFAULT_API_URL, APIKEY, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)
	authentication, _ := client.LoginWithPassword(USER, PASS)

	certificateClient := api.NewCertificateClient(apiContracts.DEFAULT_API_CERTIFICATE_URL, authentication.Token, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)

	certificateRequest := apiContracts.CertificateRequest{
		MachineID: "55eb0e08ddd14f8d8752a982e18bd4aa",
		ServiceID: "80:00:00:00:01:0C:2C:27",
		Name:      "google_com-wwws",
		IP:        "www.google.com",
	}

	certificateResponse, errx := certificateClient.Generate(certificateRequest)
	if errx != nil {
		t.Error(errx)
		t.FailNow()
	}

	if certificateResponse == nil {
		t.FailNow()
	}
}
