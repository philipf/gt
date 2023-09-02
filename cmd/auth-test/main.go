package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	auth "github.com/microsoft/kiota-authentication-azure-go"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"

	abstractions "github.com/microsoft/kiota-abstractions-go"
)

func main() {
	fmt.Println("Go Graph Tutorial")
	fmt.Println()

	// Load .env files
	// .env.local takes precedence (if present)
	err := godotenv.Load(".env.local")
	//err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env", err)
	}

	graphHelper := NewGraphHelper()

	err = initializeGraph(graphHelper)

	if err != nil {
		log.Panicf("Error initializing Graph for user auth: %v\n", err)
	}

	displayAccessToken(graphHelper)

	greetUser(graphHelper)

	var graphClient = graphHelper.userClient

	lists, err := graphClient.Me().Todo().Lists().Get(context.Background(), nil)

	if err != nil {
		fmt.Printf("Error creating lists: %v\n", err)
	}

	for _, list := range lists.GetValue() {
		headers := abstractions.NewRequestHeaders()

		var pageSize int32 = 50
		var filter = "status ne 'completed'"

		query := users.ItemTodoListsItemTasksRequestBuilderGetQueryParameters{
			Top:    &pageSize,
			Filter: &filter,
			//Select: []string{"title", "status"}, //this doesn't work
		}

		options := users.ItemTodoListsItemTasksRequestBuilderGetRequestConfiguration{
			Headers:         headers,
			QueryParameters: &query,
		}

		listId := *list.GetId()
		tasks, err := graphClient.Me().Todo().Lists().ByTodoTaskListIdString(listId).Tasks().Get(context.Background(), &options)

		if err != nil {
			log.Fatalf("Error ByTodoTaskListIdString: %v\n", err)
		}

		fmt.Printf("Tasks Data: %+v\n", tasks)

		pageIterator, err := graphcore.NewPageIterator[models.TodoTaskable](
			tasks,
			graphClient.GetAdapter(),
			// models.CreateTodoTaskFromDiscriminatorValue
			models.CreateTodoTaskCollectionResponseFromDiscriminatorValue)

		if err != nil {
			log.Fatalf("Error creating page iterator: %v\n", err)
		}

		pageIterator.SetHeaders(headers)

		var taskCount = 1

		err = pageIterator.Iterate(
			context.Background(),
			func(task models.TodoTaskable) bool {
				fmt.Printf("Task: %d] %s\n", taskCount, *task.GetTitle())
				taskCount++
				return true
			})

		if err != nil {
			log.Fatalf("Error iterating over messages: %v\n", err)
		}
	}

	//for _, list := range lists.GetValue() {
	// 	fmt.Println("List:", *list.GetDisplayName())
	// 	tasks, err := graphHelper.userClient.Me().Todo().Lists().ByTodoTaskListIdString(*list.GetId()).Tasks().Get(context.Background(), nil)

	// 	if err != nil {
	// 		fmt.Printf("Error getting tasks: %v\n", err)
	// 	}

	// 	for _, task := range tasks.GetValue() {
	// 		fmt.Println("Task:", *task.GetTitle())
	// 	}

	// 	//tasks.GetOdataNextLink()

	// 	//https://learn.microsoft.com/en-us/graph/tutorials/go?tabs=aad&tutorial-step=5
	// 	// https://learn.microsoft.com/en-us/graph/sdks/paging?tabs=go

	// 	// for t := range tasks.GetValue() {
	// 	// 	fmt.Println("Task:vv", t.GetTitle())
	// 	// }
	// }
}

type GraphHelper struct {
	deviceCodeCredential *azidentity.DeviceCodeCredential
	userClient           *msgraphsdk.GraphServiceClient
	graphUserScopes      []string
}

func NewGraphHelper() *GraphHelper {
	g := &GraphHelper{}
	return g
}

func initializeGraph(g *GraphHelper) error {
	clientId := os.Getenv("CLIENT_ID")
	tenantId := os.Getenv("TENANT_ID")
	scopes := os.Getenv("GRAPH_USER_SCOPES")
	g.graphUserScopes = strings.Split(scopes, ",")

	// Create the device code credential
	credential, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		ClientID: clientId,
		TenantID: tenantId,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Println(message.Message)
			return nil
		},
	})
	if err != nil {
		return err
	}

	g.deviceCodeCredential = credential

	// Create an auth provider using the credential
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, g.graphUserScopes)
	if err != nil {
		return err
	}

	// Create a request adapter using the auth provider
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return err
	}

	// Create a Graph client using request adapter
	client := msgraphsdk.NewGraphServiceClient(adapter)
	g.userClient = client

	return nil
}

func greetUser(graphHelper *GraphHelper) {
	user, err := graphHelper.GetUser()
	if err != nil {
		log.Panicf("Error getting user: %v\n", err)
	}

	fmt.Printf("Hello, %s!\n", *user.GetDisplayName())

	// For Work/school accounts, email is in Mail property
	// Personal accounts, email is in UserPrincipalName
	email := user.GetMail()
	if email == nil {
		email = user.GetUserPrincipalName()
	}

	fmt.Printf("Email: %s\n", *email)
	fmt.Println()
}

func (g *GraphHelper) GetUserToken() (*string, error) {
	token, err := g.deviceCodeCredential.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: g.graphUserScopes,
	})
	if err != nil {
		return nil, err
	}

	return &token.Token, nil
}

func displayAccessToken(graphHelper *GraphHelper) {
	token, err := graphHelper.GetUserToken()
	if err != nil {
		log.Panicf("Error getting user token: %v\n", err)
	}

	fmt.Printf("User token: %s", *token)
	fmt.Println()
}

func (g *GraphHelper) GetUser() (models.Userable, error) {
	query := users.UserItemRequestBuilderGetQueryParameters{
		// Only request specific properties
		Select: []string{"displayName", "mail", "userPrincipalName"},
	}

	return g.userClient.Me().Get(context.Background(),
		&users.UserItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &query,
		})
}

func PrintOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error:", typed.Error())
		if terr := typed.GetErrorEscaped(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}
