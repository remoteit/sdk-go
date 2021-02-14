package tests

import (
	"fmt"
	"testing"

	api "github.com/remoteit/sdk-go"
	apiContracts "github.com/remoteit/sdk-go/contracts"
)

func Test_GraphQL_GetServiceNamesByIDs(t *testing.T) {
	client := api.NewClient(apiContracts.DEFAULT_API_URL, APIKEY, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)
	authentication, _ := client.LoginWithPassword(USER, PASS)

	graphQLClient := api.NewGraphQLClient(apiContracts.DEFAULT_API_GRAPHQL_URL, authentication.Token, apiContracts.DEFAULT_API_TIMEOUT, apiContracts.DEFAULT_API_USER_AGENT)

	services, errx := graphQLClient.GetServiceNamesByIDs([]string{SERVICEID})
	if errx != nil {
		t.Error(errx)
		t.FailNow()
	}

	for _, service := range services {
		fmt.Println(service)
	}
}
