package tests

import (
	"fmt"
	"testing"

	api "github.com/remoteit/sdk-go"
	apiContracts "github.com/remoteit/sdk-go/contracts"
)

func Test_GraphQL_GetDeviceAndServiceNames(t *testing.T) {
	client := api.NewClient(apiContracts.DEFAULT_API_URL, APIKEY, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)
	authentication, _ := client.LoginWithPassword(USER, PASS)

	graphQLClient := api.NewGraphQLClient(apiContracts.DEFAULT_API_GRAPHQL_URL, authentication.Token, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)

	device, errx := graphQLClient.GetDeviceAndServiceNames(DEVICEID)
	if errx != nil {
		t.Error(errx)
		t.FailNow()
	}

	fmt.Println(device.ID)
	fmt.Println(device.Name)

	for _, service := range device.Services {
		fmt.Println(service)
	}
}
