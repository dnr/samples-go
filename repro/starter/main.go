package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"time"

	"go.temporal.io/sdk/client"
	"golang.org/x/time/rate"

	"github.com/temporalio/samples-go/repro"
)

func main() {
	rpsflag := flag.Float64("rps", 20, "send rps")
	gorflag := flag.Int("gor", 10, "goroutines")
	flag.Parse()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// argh :(
	go func() {
		pp, _ := os.FindProcess(os.Getppid())
		for pp.Signal(syscall.Signal(0)) == nil {
			time.Sleep(time.Second)
		}
		os.Exit(0)
	}()

	lim := rate.NewLimiter(rate.Limit(*rpsflag), 1)

	var wg sync.WaitGroup
	for w := 0; w < *gorflag; w++ {
		w := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			i := 0
			for {
				lim.Wait(context.Background())
				i++
				workflowOptions := client.StartWorkflowOptions{
					ID:        fmt.Sprintf("repro-%d-%d-%09d", os.Getpid(), w, i),
					TaskQueue: "tq",
				}
				we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, repro.Workflow, i)
				if err != nil {
					log.Fatalln("Unable to execute workflow", err)
				}
				log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
			}
		}()
	}
	wg.Wait()
}
