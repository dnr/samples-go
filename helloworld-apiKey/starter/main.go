package main

import (
	"context"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	helloworldapiKey "github.com/temporalio/samples-go/helloworld-apiKey"
	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	clientOptions, err := helloworldapiKey.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	for {
		time.Sleep(time.Second)

		workflowOptions := client.StartWorkflowOptions{
			ID:        "hello-" + strconv.Itoa(int(rand.Intn(math.MaxInt64))),
			TaskQueue: "ptest",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, helloworldapiKey.Workflow, "Temporal")
		if err != nil {
			log.Println("Unable to start workflow", err)
			continue
		}

		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	}

	// // Synchronously wait for the workflow completion.
	// var result string
	// err = we.Get(context.Background(), &result)
	// if err != nil {
	// 	log.Fatalln("Unable get workflow result", err)
	// }
	// log.Println("Workflow result:", result)
}
