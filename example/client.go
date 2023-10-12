package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/avisi-cloud/go-client/pkg/acloudapi"
)

func RunExample() {
	client := acloudapi.NewClient(
		acloudapi.NewPersonalAccessTokenAuthenticator(os.Getenv("ACLOUD_PAT")),
		acloudapi.ClientOpts{
			Debug:                        false,
			Trace:                        false,
			DebugShowAuthorizationHeader: false,
			APIUrl:                       os.Getenv("ACLOUD_URL"),
			UserAgent:                    "example",
		},
	)

	organisationSlug := "ame"
	cloudAccounts, err := client.GetCloudAccounts(context.Background(), organisationSlug)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get cloud-accounts: %v", err)
		os.Exit(1)
	}
	cloudAccountsJson, _ := json.MarshalIndent(cloudAccounts, "", "\t")
	log.Printf("received cloud-accounts:\n")
	fmt.Printf("%s\n", cloudAccountsJson)

	clusters, err := client.GetClusters(context.Background())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get clusters: %v", err)
		os.Exit(1)
	}
	clustersJson, _ := json.MarshalIndent(clusters, "", "\t")
	log.Printf("received clusters:\n")
	fmt.Printf("%s\n", clustersJson)

	organisation, err := client.GetOrganisation(context.Background(), organisationSlug)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get organisation %q: %v", organisationSlug, err)
		os.Exit(1)
	}
	organisationJson, _ := json.MarshalIndent(organisation, "", "\t")
	log.Printf("received organisation:\n")
	fmt.Printf("%s\n", organisationJson)
}
