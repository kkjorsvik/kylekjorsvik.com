---
title: "Kyvik"
description: "A security-first, multi-agent AI framework in Go — guardrails, sandboxed isolation, and a dashboard, in a single binary."
status: "active"
tags: ["go", "ai-agents", "security", "postgres"]
draft: false
featured: true
links:
  - label: "GitHub"
    url: "https://github.com/kkjorsvik/kyvik"
  - label: "Website"
    url: "https://kyvik.io/"
---

Kyvik is a personal AI agent platform I built in Go because nothing off the shelf did
exactly what I wanted. Most frameworks either hand agents unrestricted access to the
host and hope for the best, or they're powerful but require deep infrastructure
investment to run safely. I wanted something in between — a platform where I own the
runtime, the agents have real isolation, and I can extend it whenever I need to.

The primary agent running on Kyvik today is Hank — my personal ops agent. Hank handles
daily accountability check-ins over Discord, routes tasks, connects to NorthstarOS (my
task execution system), and keeps me honest about what I said I was going to do. I also
run a Job Hunter agent through Kyvik during my current job search.

The reason I built it myself instead of reaching for an existing platform is the same
reason I build most things: I learn a domain better when I'm responsible for the whole
thing. Kyvik is one of my strongest Go projects and the one I've learned the most from.

## Design goals

- **Security & guardrails first** — deny-by-default permissions on every tool call,
  sandboxed execution, spending limits, and audit logging. Built into the framework,
  not bolted on.
- **Multi-agent native** — each agent has its own identity, permissions, execution
  sandbox, and communication boundaries. Isolation by design.
- **Go-native simplicity** — single-binary deployment, low resource footprint, no
  runtime dependencies beyond the binary and PostgreSQL.
- **Extensible** — when I need a new capability, I add it. That's the point of owning
  the platform.

## Status

Active development. Core runtime is implemented — agent lifecycle management, permission
gates, model routing across OpenRouter, OpenAI, Anthropic, and Ollama, Discord
integration, and the web dashboard. Open source under the MIT license with a potential
SaaS path down the road.
