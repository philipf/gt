package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

var (
	tenantdID = "1234567890"
	clientID  = "1234567890"

// clientSecret = "1234567890"
)

func main() {
	ctx := context.Background()
	env := environments.AzurePublic()

	credentials := auth.Credentials{
		Environment: *env,
		TenantID:    tenantdID,
		ClientID:    clientID,
		//ClientSecret: clientSecret,
		EnableAuthenticatingUsingAzureCLI: true,
	}

	authorizer, err := auth.NewAuthorizerFromCredentials(ctx, credentials, env.MicrosoftGraph)
	if err != nil {
		log.Fatal(err)
	}

	client := msgraph.NewMeClient()
	client.BaseClient.Authorizer = authorizer

	me, _, err := client.Get(ctx, odata.Query{})

	if err != nil {
		log.Println(err)
		return
	}
	if me == nil {
		log.Println("bad API response, nil result received")
		return
	}

	fmt.Printf("Hello, %s!\n", *me.DisplayName)
}
