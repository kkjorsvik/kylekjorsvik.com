---
title: "Quarterdeck"
description: "A Linux-native IDE for the AI agent era — one keyboard-driven control plane for AI coding agents across projects."
status: "paused"
tags: ["go", "wails", "react", "typescript"]
draft: false
links:
  - label: "GitHub"
    url: "https://github.com/kkjorsvik/quarterdeck"
---

A Linux-native IDE for the AI agent era.

Quarterdeck is a unified workspace for developers who work across multiple projects at once by
delegating to AI coding agents. Instead of juggling terminals, editors, and git tools scattered
across your desktop, it brings the entire agent workflow — spawning, monitoring, reviewing, and
merging — into one keyboard-driven environment.

You can launch and track AI coding agents (Claude Code, Codex, OpenCode) across several
codebases simultaneously, watch their live terminal output, review their changes per-run with a
built-in diff viewer, and accept, reject, or commit their work — all from a single window.

## What it does

- **Multi-project workspace** — manage many codebases at once with per-project layout memory and instant switching
- **Agent management** — spawn, monitor, and track agents with state detection and desktop notifications; re-run agents with stored parameters; per-project health summaries
- **Real terminals** — full PTY-based terminal emulation (xterm.js + WebGL) that persists in the background across project switches, with output saved to disk for post-mortem review
- **Diff & review** — review each agent run side-by-side in a Monaco diff editor, accept/reject per file, and commit with pre-populated messages
- **Git integration** — worktree isolation per agent, branch management, conflict resolution, status indicators, log, and stash
- **Built-in editor** — Monaco with vim keybindings, syntax highlighting, and tabbed editing
- **i3-inspired tiling** — keyboard-driven panels with splits, tabs, and focus management

## Design goals

- **Linux-native, not Electron** — a compiled Go binary (~45 MB) using WebKitGTK, built for minimal resource usage
- **Keyboard-first** — an i3-inspired tiling layout for developers who live in a tiling window manager
- **One control plane for many agents** — treat AI agents as first-class, parallel workers across projects rather than one-at-a-time CLI sessions
- **Built for a real workflow** — a personal tool shaped around working across multiple projects on a tiling-WM Linux setup

## Tech stack

Go + Wails v2 backend · React 18 + TypeScript + Zustand frontend · xterm.js terminals · Monaco
editor · SQLite (pure-Go, no CGO).

## Status

Paused. Development is currently on hold. The project reached a working state — including a
significant refactor from a general IDE toward a dedicated agent control plane — but I'm
reconsidering its core architecture. The terminal-emulation approach (wrapping agent CLIs in
PTY + WebSocket + xterm.js and detecting agent state from their output) proved heavyweight and
fragile. Any future direction would likely move toward talking to agents through their APIs/SDKs
for structured data, à la Zed — a fundamentally different architecture worth getting right
before building further.
