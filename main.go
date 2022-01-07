package main

import (
	"context"
	"fmt"
	"time"
)

const defaultTimeout = time.Second * 2
const apiCallExecutionTime = time.Second * 4

func main() {
	fmt.Println("Setup: API call takes " + apiCallExecutionTime.String() + " to execute")

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	fmt.Println("\n...executing API call with default timeout (" + defaultTimeout.String() + ") context...")
	fmt.Println(execute(ctx))

	ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	fmt.Println("\n...executing API call with default timeout (" + defaultTimeout.String() + ") but detached context...")
	fmt.Println(executeWithDetachedContext(ctx))

	ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	extendedTimeout := time.Second * 3
	fmt.Println("\n...executing API call with extended timeout (" + extendedTimeout.String() + ") and detached context...")
	fmt.Println(executeWithDetachedAndExtendedContext(ctx, extendedTimeout))

	ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	extendedTimeout = time.Second * 5
	fmt.Println("\n...executing API call with extended timeout (" + extendedTimeout.String() + ") and detached context...")
	fmt.Println(executeWithDetachedAndExtendedContext(ctx, extendedTimeout))
}

func execute(ctx context.Context) string {
	return someApiCall(ctx)
}

func executeWithDetachedContext(ctx context.Context) string {
	detachedCtx := DetachedContext(ctx)
	return someApiCall(detachedCtx)
}

func executeWithDetachedAndExtendedContext(ctx context.Context, timeout time.Duration) string {
	detachedAndExtendedCtx, cancel := context.WithTimeout(DetachedContext(ctx), timeout)
	defer cancel()
	return someApiCall(detachedAndExtendedCtx)
}

func someApiCall(ctx context.Context) string {
	ch := make(chan string)

	go func() {
		time.Sleep(apiCallExecutionTime)
		ch <- "API call response: OK"
	}()

	select {
	case <-ctx.Done():
		return ctx.Err().Error()
	case res := <-ch:
		return res
	}
}

// DetachedContext returns a new context detached from the lifetime
// of ctx, but which still returns the values of ctx.
//
// DetachedContext can be used to maintain the trace context required
// to correlate events, but where the operation is "fire-and-forget",
// and should not be affected by the deadline or cancellation of ctx.

type detachedContext struct {
	context.Context
	orig context.Context
}

func DetachedContext(ctx context.Context) context.Context {
	return &detachedContext{Context: context.Background(), orig: ctx}
}
