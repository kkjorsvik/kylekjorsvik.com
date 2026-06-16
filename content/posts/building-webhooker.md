---
title: "Building a Webhook Delivery Platform in Go"
date: "2026-05-28"
tags: ["go", "kubernetes", "webhooks"]
draft: true
description: "How I'm building Webhooker — a Go monorepo for reliable webhook delivery with Postgres, Redis, and asynq."
---

Webhook delivery sounds trivial until you actually have to make it reliable. Retries,
idempotency, signature verification, dead-letter queues, per-tenant rate limiting — the
"just POST a JSON body" mental model falls apart fast. Webhooker is my attempt to build
that machinery properly, as a deliberate learning vehicle.

## Why Webhooker

I wanted a project that forced me to touch the parts of distributed systems I usually only
configure rather than build: durable queues, exactly-enough-once semantics, and graceful
degradation under load.

The core delivery loop is a worker pulling jobs off a queue:

```go
package main

import (
	"context"
	"log"
)

func deliver(ctx context.Context, job DeliveryJob) error {
	resp, err := client.Post(ctx, job.Endpoint, job.Payload)
	if err != nil {
		return retryable(err)
	}
	if resp.StatusCode >= 500 {
		return ErrRetry
	}
	return nil
}

func main() {
	log.Println("webhooker worker started")
}
```

## Architecture

- **Postgres** for durable state — endpoints, delivery attempts, audit log
- **Redis + asynq** for the work queue and exponential-backoff retries
- **A thin HTTP ingress** that validates, persists, and enqueues

More to come as the project matures.
