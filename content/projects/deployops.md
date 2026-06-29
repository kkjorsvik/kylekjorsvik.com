---
title: "DeployOps"
description: "An audit-ready deployment operations platform for regulated teams — tracking, approvals, drift detection, and release management."
status: "active"
tags: ["go", "react", "postgres", "compliance"]
draft: false
featured: true
links:
  - label: "Website"
    url: "https://deployops.io/"
  - label: "Blog post"
    url: "/blog/building-deployops/"
---

A deployment operations platform for regulated teams.

DeployOps gives compliance-conscious engineering organizations an audit-ready system of record
for everything that ships to production. It replaces fragile spreadsheets and scattered tooling
with a single platform for deployment tracking, multi-stakeholder approvals, configuration
drift detection, and structured release management.

The core idea: when an auditor asks "Show me every production deployment in the last six months
— who approved each one and what changed," the answer takes five minutes, not five days. It's
built for the 50–500 person company in a regulated space (healthcare SaaS, fintech, government
contractors) navigating SOC 2 and similar frameworks.

## What it does

- **Deployment tracking** — record every release via REST API or web UI, integrating with Jenkins, GitHub Actions, or any CI/CD pipeline
- **Approval workflows** — multi-stakeholder sign-off with email notifications, secure token-based approval links, and CI/CD gating
- **State verification** — SSH-based scanning that detects configuration drift across environments
- **Release cycle management** — version tracking, multi-service deployments, and a service × environment status matrix
- **Dynamic release notes** — customizable, drag-and-drop release documentation with JSON export for downstream automation
- **Audit reporting** — one-click PDF reports with complete approval trails
- **SDLC path enforcement** — tracks promotions through DEV → SIT → UAT → PROD with warnings for out-of-order changes
- **Security** — organization-scoped multi-tenancy, role-based access control, and TOTP two-factor authentication

## Tech stack

- **Backend** — Go 1.25, chi router, PostgreSQL 16 (pgx), Prometheus metrics, structured slog logging.
- **Frontend** — React 18, TypeScript, Vite, Tailwind CSS, with API client types generated from an OpenAPI 3 contract.
- **Infrastructure** — multi-stage Docker builds, GitHub Actions CI, deployable via Docker Compose or a one-line install script.

## Design goals

DeployOps is built around a deliberate engineering philosophy and a multi-phase hardening
roadmap:

- **Compliance-first** — immutable, append-only audit logs; every action attributable and exportable
- **API-first** — every feature is available over REST, with the OpenAPI contract as the single source of truth and generated client types that keep frontend and backend in lockstep
- **Observable by default** — request-ID correlation on every response, structured JSON logs, and Prometheus metrics from day one
- **Correct under failure** — transaction-wrapped mutations, errors that map cleanly to HTTP responses without leaking internals, and rollback-tested migrations
- **Quality that doesn't regress** — CI quality gates, linting and tests, and a coverage ratchet to keep standards from slipping

## Status

Active development. The platform's seven feature areas — from deployment tracking through
multi-service release cycles — are implemented and in use, and work is currently focused on the
hardening roadmap (observability, API contracts, transactional integrity, and error-handling
consistency) toward production-grade reliability.
