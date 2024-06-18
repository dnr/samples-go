package fairpri

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context, tenantId string) (string, error) {
	l := workflow.GetLogger(ctx)
	l.Info("Workflow started", "tenant", tenantId)

	actCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		FairnessKey:         tenantId, // Defaults to workflow's fairness key but can be overridden
	})
	err := workflow.ExecuteActivity(actCtx, ActivityOne).Get(ctx, nil)
	if err != nil {
		l.Error("ActivityOne failed", "Error", err)
		return "", err
	}

	// Can read workflow's priority out of info:
	priority := workflow.GetInfo(ctx).Priority
	if wereInAHurry {
		priority += 100
	}

	actCtx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		Priority:            priority, // Override priority
	})
	err := workflow.ExecuteActivity(actCtx, ActivityTwo).Get(ctx, nil)
	if err != nil {
		l.Error("ActivityTwo failed", "Error", err)
		return "", err
	}

	return result, nil
}

func ActivityOne(ctx context.Context) error {
	return nil
}

func ActivityTwo(ctx context.Context) error {
	// Can read activity's priority out of info, e.g. to propagate to downstream system,
	// set nice level, etc.
	pri := activity.GetInfo(ctx).Priority

	return nil
}
