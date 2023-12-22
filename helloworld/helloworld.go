package helloworld

import (
	"fmt"
	"sync/atomic"
	"time"

	"go.temporal.io/sdk/workflow"

	// TODO(cretz): Remove when tagged
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

var (
	ctr1 atomic.Int32
	ctr4 atomic.Int32
)

func Workflow(ctx workflow.Context, name string) (string, error) {
	fmt.Printf("@@@ STARTING WORKFLOW\n")

	for i := 0; i < 5; i++ {
		fmt.Printf("@@@ i == %d\n", i)
		nowBefore := workflow.Now(ctx)

		var nowFromSideEffect time.Time
		workflow.SideEffect(ctx, func(ctx workflow.Context) any {
			fmt.Printf("@@@ in SideEffect, returning %v\n", nowBefore)
			return nowBefore
		}).Get(&nowFromSideEffect)

		// fmt.Printf("@@@ nowBefore %v, nowFromSideEffect %v\n", nowBefore, nowFromSideEffect)
		if nowBefore != nowFromSideEffect {
			fmt.Printf("@@@ !!!!!!!!!!!!!!!!!!!!!! different\n")
			panic(false)
		}
		if nowBefore != workflow.Now(ctx) {
			fmt.Printf("@@@ !!!!!!!!!!!!!!!!!!!!!! different2\n")
			panic(false)
		}

		switch i {
		case 1:
			switch ctr1.Add(1) {
			case 1, 2, 3:
				fmt.Printf("@@@ forcing timeout, i %d, ctr %d\n", i, ctr1.Load())
				time.Sleep(workflow.GetInfo(ctx).WorkflowTaskTimeout - 500*time.Millisecond)
			}
		case 4:
			switch ctr4.Add(1) {
			case 1:
				fmt.Printf("@@@ forcing timeout, i %d, ctr %d\n", i, ctr4.Load())
				time.Sleep(workflow.GetInfo(ctx).WorkflowTaskTimeout - 500*time.Millisecond)
			}
		}

		if nowBefore != workflow.Now(ctx) {
			fmt.Printf("@@@ !!!!!!!!!!!!!!!!!!!!!! different3\n")
			panic(false)
		}

		fmt.Printf("@@@ wf sleeping\n")
		workflow.Sleep(ctx, 2*time.Second)
	}

	fmt.Printf("@@@ wf done\n")
	return "ok", nil
}
