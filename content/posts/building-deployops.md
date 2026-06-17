---
title: "Building DeployOps: Visibility Into What's Actually Running in Production"
date: "2026-06-17"
tags: ["go", "devops", "deployops", "compliance"]
draft: true
description: "I built DeployOps because I couldn't answer a basic question about my own production environment without SSHing into boxes. Here's what it does and how it works."
---

When an auditor asks "show me every production deployment from the last six months —
who approved each one and what changed" — how long does that take your team to answer?

For a lot of engineering organizations in regulated spaces, the honest answer is days.
Maybe longer. The actual record is scattered across Slack threads, Confluence pages that
haven't been updated since Q2, and a spreadsheet someone made during the last audit that
nobody maintained after it was over.

I built DeployOps because I hit a simpler version of that problem first: I couldn't
tell you what version of each microservice was actually running in production without
logging into every server individually and reading the Docker Compose file. That's not
an audit problem. That's just a visibility problem. And it bothered me enough to build
something about it.

## The problem that started it

I work in a regulated environment — government healthcare SaaS, NIST 800-171. We have
a lot of services. When I want to know what's deployed where, my options are: ask
whoever did the last deployment and hope they remember, dig through Jenkins build logs,
or SSH into each host and check manually.

None of those scale. None of them give you drift detection. And none of them tell you
if someone went into a server and changed something outside of the normal deployment
process — the "wild west" change that bypasses every approval and leaves no trace in
your CI/CD logs.

I wanted a system where every deployment reports in, every scan of the environment
reports in, and anything that doesn't match what's expected gets flagged immediately.
I haven't gotten clearance to run it at work yet — politics and process move slower
than code — but I built it anyway because the problem is real and the tool should exist.

## What DeployOps actually does

DeployOps is a system of record and a gate. It does not run deployments. That
distinction matters.

Your existing pipelines — Jenkins, GitHub Actions, whatever you're using — stay exactly
as they are. You add one API call at the end of each deployment to tell DeployOps what
just changed. A scan pipeline runs on a schedule, SSHs into your hosts, reads the
running state of your services, and reports what it finds back to DeployOps via API.
DeployOps compares what it finds against what it expects and flags anything that doesn't
match.

That's the core loop. Everything else — approval workflows, release notes, audit
reports, SDLC path enforcement — is built on top of that foundation.

## How drift detection works

The first scan establishes the baseline. Every service, every version, every
environment. That becomes the expected state.

From that point forward, the only thing that can legitimately change the expected state
is a deployment. When a deployment comes in via API, DeployOps updates its expectations
for the services that were deployed. Everything else stays where it was.

When the next scan runs, it compares what it finds against what's expected. Two
scenarios trigger a drift alert:

**Post-deployment drift** — a service was just deployed but the scan finds a different
version than what was reported in the deployment. Something went wrong between the
deployment reporting in and the service actually running the new version.

**Wild west drift** — no deployment happened, but the scan finds a version that
doesn't match the last known good state. Someone changed something on the server
directly. Maybe it was intentional. Maybe it was an emergency fix. Either way, it
didn't go through any process, it's not tracked anywhere, and now DeployOps knows
about it.

That second scenario is the one that matters most in a regulated environment. It's the
change that doesn't show up in your deployment logs, doesn't have an approver, and
would be invisible without active scanning. DeployOps makes it visible.

## The approval workflow

For teams that need it, DeployOps has a multi-stakeholder approval process built in.
Deployments can require sign-off before pipelines are allowed to proceed — pipelines
check in with DeployOps before running and get a go/no-go response based on whether
the required approvals have been collected.

Approvals are done via secure token-based links sent by email. No account required for
approvers. The full approval trail — who approved, when, from where — is stored as an
immutable append-only record. When an auditor asks for it, you export a PDF. Five
minutes, not five days.

The system also enforces SDLC path ordering. If a change tries to promote to UAT
before it's cleared SIT, DeployOps flags it. Out-of-order promotions are a common
compliance finding and a surprisingly easy thing to accidentally do when you're moving
fast.

## Technical decisions worth mentioning

A few things I'm reasonably happy with on the implementation side:

**API-first with generated types** — the OpenAPI spec is the source of truth. The
React frontend uses types generated directly from that spec, so the frontend and backend
can't drift out of sync without it being caught at build time. This was worth the
upfront investment.

**Append-only audit records** — nothing in the audit trail gets updated or deleted.
Every action appends a new record. This is the right design for compliance use cases
and it simplifies a lot of edge cases around concurrent writes.

**Real Postgres test harness** — the test suite runs against an actual Postgres
instance, not mocks. This caught several bugs that a mocked database would have missed,
particularly around transaction behavior and constraint violations. It's slower to set
up but I won't go back to mocking the database layer.

**Observable by default** — every log line and error response carries a request ID.
Prometheus metrics from day one. When something goes wrong in production you want to
be able to trace a request through the system without adding instrumentation after the
fact.

## Current state

The core feature set is implemented and working — deployment tracking, drift detection,
approval workflows, release notes, the service × environment status matrix, audit
reporting. The work right now is hardening: transactional integrity, consistent error
handling, API contract enforcement, and test coverage that doesn't regress.

The goal is a self-hostable platform that a 50–500 person engineering team in a
regulated space can run on their own infrastructure with minimal ops overhead. Docker
Compose or a one-line install script. No vendor lock-in, no SaaS dependency for
something as critical as your deployment audit trail.

I plan to open source it once the hardening phase is done and I'm happy with the
quality of what I'm putting out. Whether it eventually becomes a hosted SaaS offering
is something I'm still thinking through — the use case fits the model, but I want the
self-hosted version to be genuinely good first.

If you're running services in a regulated environment and you recognize the problems
described here, I'd be interested to hear how you're solving them today.
