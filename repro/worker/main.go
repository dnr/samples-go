package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/repro"
)

func main() {
	pollersflag := flag.Int("pollers", 2, "wft pollers")
	tpsflag := flag.Float64("tps", 35, "local activities per second (rate limit wfts)")
	flag.Parse()

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "tq", worker.Options{
		MaxConcurrentWorkflowTaskPollers: *pollersflag,
		// MaxConcurrentWorkflowTaskExecutionSize: 20,
		WorkerLocalActivitiesPerSecond: *tpsflag,
	})

	w.RegisterWorkflow(repro.Workflow)
	w.RegisterActivity(repro.LocalAct)

	// argh :(
	go func() {
		pp, _ := os.FindProcess(os.Getppid())
		for pp.Signal(syscall.Signal(0)) == nil {
			time.Sleep(time.Second)
		}
		w.Stop()
	}()

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
