package helloworldapiKey

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a Hello World workflow definition.
func Workflow(ctx workflow.Context, name string) error {
	const acts = 10
	fkeys := []string{"A", "B", "C"}

	logger := workflow.GetLogger(ctx)

	var futs []workflow.Future
	for _, fkey := range fkeys {
		for range acts {
			ao := workflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
				TaskQueue:           "fairtest2",
				Priority:            temporal.Priority{FairnessKey: fkey},
			}
			ctx = workflow.WithActivityOptions(ctx, ao)
			fut := workflow.ExecuteActivity(ctx, Activity)
			futs = append(futs, fut)
		}
	}

	for _, f := range futs {
		err := f.Get(ctx, nil)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return err
		}
	}

	return nil
}

func Activity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	fkey := activity.GetInfo(ctx).Priority.FairnessKey
	logger.Info("Activity", "fkey", fkey)
	return nil
}

// ParseClientOptionFlags parses the given arguments into client options. In
// some cases a failure will be returned as an error, in others the process may
// exit with help info.
func ParseClientOptionFlags(args []string) (client.Options, error) {
	// Parse args
	set := flag.NewFlagSet("hello-world-api-key", flag.ExitOnError)
	targetHost := set.String("target-host", "us-east-1.aws.api.temporal.io:7233", "Host:port for the server")
	namespace := set.String("namespace", "david.a2dd6", "Namespace for the server")
	apiKey := set.String("api-key", "", "Optional API key, mutually exclusive with cert/key")

	if err := set.Parse(args); err != nil {
		return client.Options{}, fmt.Errorf("failed parsing args: %w", err)
	}

	if *apiKey == "" {
		*apiKey = os.Getenv("TEMPORAL_CLIENT_API_KEY")
	}
	if *apiKey == "" {
		return client.Options{}, fmt.Errorf("-api-key or TEMPORAL_CLIENT_API_KEY env is required required")
	}

	return client.Options{
		HostPort:          *targetHost,
		Namespace:         *namespace,
		ConnectionOptions: client.ConnectionOptions{TLS: &tls.Config{}},
		Credentials:       client.NewAPIKeyStaticCredentials(*apiKey),
	}, nil
}
