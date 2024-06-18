package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/fairpri"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// From client's pov, it listens on a single "task queue".
	w := worker.New(c, "tq", worker.Options{
		// If true, do an extra poll for activities with higher priority than what we have
		// running. If we get one, cancel current lowest-priority activity and run that
		// instead. (Phase 2+)
		PreemptActivities: true,
	})

	w.RegisterWorkflow(fairpri.Workflow)
	w.RegisterActivity(fairpri.Activity)

	if err = w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
