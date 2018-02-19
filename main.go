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

	"os"

	"strconv"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v0.beta"
	"gopkg.in/mailgun/mailgun-go.v1"
	"net/http"
)

func main() {

	http.HandleFunc("/", nullHandler)
	http.HandleFunc("/check_quotas", notifyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func nullHandler(w http.ResponseWriter, r *http.Request) {

}

func logHandler(w http.ResponseWriter, r *http.Request) {
	quotas := getQuotasToLog()
	for _, quota := range quotas {
		fmt.Printf("EXCEEDED QUOTA: %s", quota)
	}
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	quotas := getQuotasToLog()
	quotaString := ""
	for _, quota := range quotas {
		quotaString = fmt.Sprintf("%s\n", quota)
	}
	mg := mailgun.NewMailgun(os.Getenv("MG_DOMAIN"), os.Getenv("MG_API_KEY"), os.Getenv("MG_PUBLIC_API_KEY"))
	message := mg.NewMessage(
		os.Getenv("MG_FROM_EMAIL"),
		"GCP Quotas Above Utilization Threshold",
		quotaString,
		os.Getenv("MG_TO_EMAIL"))
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}

func getQuotasToLog() []string {

	quotasToLog := []string{}
	project := os.Getenv("PROJECT_ID")
	threshold_string := os.Getenv("THRESHOLD")
	threshold, err := strconv.ParseFloat(threshold_string, 64)

	if err != nil {
		log.Fatal(err)
	}

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
			quotasToLog = append(quotasToLog, fmt.Sprintf("%s - %#v\n", quota.Metric, utilized))
		}

	}

	regions := []string{"us-east1", "us-east4", "us-central1", "us-west1"} // TODO: Update with the regions you're using in this project

	for _, region := range regions {
		regionResp, err := computeService.Regions.Get(project, region).Context(ctx).Do()
		if err != nil {
			log.Fatal(err)
		}

		for _, regquota := range regionResp.Quotas {
			utilized := regquota.Usage / regquota.Limit
			if utilized > threshold {
				quotasToLog = append(quotasToLog, fmt.Sprintf("%s | %s - %#v\n", region, regquota.Metric, utilized))
			}
		}
	}
	return quotasToLog
}
