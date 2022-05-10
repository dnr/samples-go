package helloworld

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type params struct {
	Val1, Val2 int
}

var defaultParams = params{12, 19}

// Workflow is a Hello World workflow definition.
func Workflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	// load default params
	get := func(ctx workflow.Context) interface{} { return defaultParams }
	eq := func(a, b interface{}) bool { return a.(params) == b.(params) }
	var params params
	if err := workflow.MutableSideEffect(ctx, "defaultParams", get, eq).Get(&params); err != nil {
		panic("can't decode params:" + err.Error())
	}

	for i := 0; i < 100; i++ {
		logger.Info("look at my params", "val1", params.Val1, "val2", params.Val2)
		workflow.Sleep(ctx, 3*time.Second)
	}

	return "ok", nil
}
