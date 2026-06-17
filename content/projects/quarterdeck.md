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

Quarterdeck is a Linux-native desktop application built with Wails v2 — Go backend,
React/TypeScript frontend, WebKitGTK rendering. The idea was a keyboard-driven control
plane for AI coding agents: PTY terminal emulation, multi-project workspace management,
agent lifecycle management for Claude Code and similar tools, a diff review workflow,
git and worktree integration, and SQLite persistence. All of it in a single native
binary without Electron.

The core problem it was solving is real — jumping between terminal windows, editor
instances, and agent outputs across multiple projects is genuinely annoying. Quarterdeck
was meant to be the thing that held it all together.

It's paused for now. The tooling landscape moved quickly while I was building it and
other people's solutions got good fast. I also had other projects pulling my attention
toward things I was enjoying more. The approach would need a rethink before I'd pick it
back up. It might come back, it might not.

## Tech stack

- **Backend** — Go, Wails v2, SQLite
- **Frontend** — React, TypeScript
- **Rendering** — WebKitGTK (not Electron)
