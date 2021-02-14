package tests

import (
	"fmt"
	"testing"

	api "github.com/remoteit/sdk-go"
	apiContracts "github.com/remoteit/sdk-go/contracts"
)

func Test_GraphQL_GetApplicationType(t *testing.T) {
	client := api.NewClient(apiContracts.DEFAULT_API_URL, APIKEY, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)
	authentication, _ := client.LoginWithPassword(USER, PASS)

	graphQLClient := api.NewGraphQLClient(apiContracts.DEFAULT_API_GRAPHQL_URL, authentication.Token, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)

	serviceType, errx := graphQLClient.GetApplicationType(SERVICEID)
	if errx != nil {
		t.Error(errx)
		t.FailNow()
	}

	if serviceType <= 0 {
		t.Error("serviceType")
		t.FailNow()
	}

	fmt.Println(serviceType)
}
