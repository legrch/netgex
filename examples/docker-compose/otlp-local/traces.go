package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// GenerateSimpleTrace creates a simple trace with a few spans
func GenerateSimpleTrace() {
	ctx, span := tracer.Start(context.Background(), "SimpleOperation")
	defer span.End()

	// Add some attributes to the span
	span.SetAttributes(
		attribute.String("operation.type", "simple"),
		attribute.Int("operation.count", 1),
	)

	// Simulate work
	time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)

	// Create a child span
	ctx, childSpan := tracer.Start(ctx, "SimpleDatabase.Query")
	childSpan.SetAttributes(attribute.String("db.operation", "select"))
	time.Sleep(time.Duration(25+rand.Intn(75)) * time.Millisecond)
	childSpan.End()
}

// GenerateComplexTrace creates a more complex trace with multiple spans and levels
func GenerateComplexTrace() {
	ctx, rootSpan := tracer.Start(context.Background(), "ComplexOperation")
	defer rootSpan.End()

	// Add some attributes to the root span
	rootSpan.SetAttributes(
		attribute.String("operation.type", "complex"),
		attribute.Int("operation.count", rand.Intn(10)+1),
		attribute.String("operation.id", fmt.Sprintf("op-%d", rand.Intn(1000))),
	)

	// Simulate some initial work
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	// Create first level of child spans with concurrency
	var wg sync.WaitGroup
	numConcurrentOperations := rand.Intn(3) + 2 // 2-4 concurrent operations

	for i := 0; i < numConcurrentOperations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			processWorkItem(ctx, index)
		}(i)
	}

	// Wait for concurrent operations to finish
	wg.Wait()

	// Finalize the trace with a closing operation
	ctx, finalizeSpan := tracer.Start(ctx, "Finalize")
	time.Sleep(time.Duration(25+rand.Intn(50)) * time.Millisecond)
	finalizeSpan.End()
}

// processWorkItem processes a single work item as part of a complex trace
func processWorkItem(ctx context.Context, index int) {
	// Create a span for this work item
	ctx, processSpan := tracer.Start(ctx, fmt.Sprintf("ProcessItem-%d", index))
	defer processSpan.End()

	processSpan.SetAttributes(
		attribute.Int("item.index", index),
		attribute.String("item.processor", fmt.Sprintf("worker-%d", index)),
	)

	// Simulate varying work duration
	time.Sleep(time.Duration(75+rand.Intn(150)) * time.Millisecond)

	// Randomly add an error to some spans
	if rand.Intn(5) == 0 {
		processSpan.SetAttributes(attribute.String("error", "process_failure"))
	}

	// Create child spans for database and API operations
	createNestedSpans(ctx, index)
}

// createNestedSpans creates database and API call child spans
func createNestedSpans(ctx context.Context, index int) {
	// Database operation
	_, dbSpan := tracer.Start(ctx, "Database.Query")
	dbSpan.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "select"),
		attribute.String("db.statement", fmt.Sprintf("SELECT * FROM items WHERE id = %d", index)),
	)
	time.Sleep(time.Duration(30+rand.Intn(70)) * time.Millisecond)
	dbSpan.End()

	// External API call
	apiCtx, apiSpan := tracer.Start(ctx, "ExternalAPI.Call")
	apiSpan.SetAttributes(
		attribute.String("http.method", "GET"),
		attribute.String("http.url", fmt.Sprintf("https://api.example.com/items/%d", index)),
	)
	time.Sleep(time.Duration(40+rand.Intn(100)) * time.Millisecond)

	// Sometimes create a nested API call
	if rand.Intn(3) == 0 {
		_, nestedSpan := tracer.Start(apiCtx, "ExternalAPI.NestedCall")
		nestedSpan.SetAttributes(
			attribute.String("http.method", "GET"),
			attribute.String("http.url", "https://api.example.com/details"),
		)
		time.Sleep(time.Duration(20+rand.Intn(50)) * time.Millisecond)
		nestedSpan.End()
	}

	apiSpan.End()
}
