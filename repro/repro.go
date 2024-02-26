package repro

import (
	"context"
	"time"

	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context, i int) (int, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("workflow started", "i", i)
	var out string
	workflow.ExecuteLocalActivity(workflow.WithLocalActivityOptions(ctx, workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: 5 * time.Second,
	}), LocalAct).Get(ctx, &out)
	return i, nil
}

func LocalAct(ctx context.Context) (string, error) {
	time.Sleep(10 * time.Millisecond)
	return "ok", nil
}
