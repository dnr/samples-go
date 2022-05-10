package helloworld

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

var defaultNumber = 19

// Workflow is a Hello World workflow definition.
func Workflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	// load default number
	get := func(ctx workflow.Context) interface{} { return defaultNumber }
	eq := func(a, b interface{}) bool { return a.(int) == b.(int) }
	var number int
	if err := workflow.MutableSideEffect(ctx, "defaultNumber", get, eq).Get(&number); err != nil {
		panic("can't decode number:" + err.Error())
	}

	for i := 0; i < 100; i++ {
		logger.Info("look at my number", "number", number)
		workflow.Sleep(ctx, 3*time.Second)
	}

	return "ok", nil
}
