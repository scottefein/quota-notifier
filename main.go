package main

// BEFORE RUNNING:
// ---------------
// 1. If not already done, enable the Compute Engine API
//    and check the quota for your project at
//    https://console.developers.google.com/apis/api/compute
// 2. This sample uses Application Default Credentials for authentication.
//    If not already done, install the gcloud CLI from
//    https://cloud.google.com/sdk/ and run
//    `gcloud beta auth application-default login`.
//    For more information, see
//    https://developers.google.com/identity/protocols/application-default-credentials
// 3. Install and update the Go dependencies by running `go get -u` in the
//    project directory.

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v0.beta"
)

func main() {

	// Project ID for this request.
	project := "scott-demo-project" // TODO: Update placeholder value.
	threshold := .2                 //What threshold you want to increase your quotas at

	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	projResp, err := computeService.Projects.Get(project).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	for _, quota := range projResp.Quotas {
		utilized := quota.Usage / quota.Limit
		if utilized > threshold {
			fmt.Printf("%s - %#v\n", quota.Metric, utilized)
		}

	}

	regions := []string{"us-east1", "us-central1", "us-west1"} // TODO: Update with the regions you're using in this project

	for _, region := range regions {
		regionResp, err := computeService.Regions.Get(project, region).Context(ctx).Do()
		if err != nil {
			log.Fatal(err)
		}

		for _, regquota := range regionResp.Quotas {
			utilized := regquota.Usage / regquota.Limit
			if utilized > threshold {
				fmt.Printf("%s - %#v\n", regquota.Metric, utilized)
			}
		}
	}

}
