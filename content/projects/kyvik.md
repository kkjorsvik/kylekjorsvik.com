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

Kyvik is a security-first, multi-agent AI framework written in Go. It provides a managed
environment for running AI agents with built-in guardrails, native multi-agent isolation, and a
web dashboard for non-technical users — all deployable as a single binary with PostgreSQL for
storage.

Most agent frameworks force a bad trade-off. One camp hands agents unrestricted access to the
host — file system, shell, credentials — and hopes nothing goes wrong. The other (LangGraph,
CrewAI, AutoGen) is powerful but demands deep Python and infrastructure expertise, treating
security as the implementer's problem. Kyvik closes that gap: it makes security boundaries the
foundational design principle while staying approachable enough that creating an agent feels
like filling out a form.

## Design goals

- **Security & guardrails first** — Deny-by-default permissions on every tool call, sandboxed execution in isolated child processes, spending limits, and audit logging for every action. These are built into the framework, not bolted on.
- **Accessible** — A web dashboard from day one, with secure defaults that protect users who don't know what's dangerous.
- **Multi-agent native** — Each agent has its own identity, permissions, execution sandbox, and communication boundaries. Agents are isolated by design, not by workaround.
- **Go-native simplicity** — Single-binary deployment, low resource footprint, goroutine-based concurrency, and no runtime dependencies beyond the binary itself.

## Status

Active development. Core runtime is implemented and MVP features are being completed — agent
lifecycle management, permission gates, sandboxing, model routing across multiple LLM providers
(OpenRouter, OpenAI, Anthropic, Ollama), and the web dashboard. Open source under the MIT
license.
