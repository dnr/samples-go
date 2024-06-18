package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/samples-go/fairpri"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	ctx := context.Background()

	tenantId := fmt.Sprintf("tenant-%d", rand.Int63())

	opts := client.StartWorkflowOptions{
		ID:          "do-stuff-for-" + tenantId + "-" + time.Now().Format(time.RFC3339),
		TaskQueue:   "tq",     // From client's pov, we're using a single "task queue"
		FairnessKey: tenantId, // FairnessKey is arbitrary string. Max length 100 bytes. Default ""
		Priority:    10,       // Priority is arbitrary int32. Higher runs first. Default 0
	}
	we, err := c.ExecuteWorkflow(ctx, opts, fairpri.Workflow, tenantId)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result string
	if err = we.Get(ctx, &result); err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
