package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/avisi-cloud/go-client/pkg/acloudapi"
)

func RunAdminExample() {
	client := acloudapi.NewAdminClient(
		acloudapi.NewPersonalAccessTokenAuthenticator(os.Getenv("ACLOUD_PAT")),
		acloudapi.ClientOpts{
			Debug:                        false,
			Trace:                        false,
			DebugShowAuthorizationHeader: false,
			APIUrl:                       os.Getenv("ACLOUD_URL"),
			UserAgent:                    "example",
		},
	)

	organisationIdentity := "ame"
	organisation, err := client.GetOrganisation(context.Background(), organisationIdentity)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get organisation: %v", err)
		os.Exit(1)
	}
	organisationJson, _ := json.MarshalIndent(organisation, "", "\t")
	log.Printf("received organisation:\n")
	fmt.Printf("%s\n", organisationJson)
}
