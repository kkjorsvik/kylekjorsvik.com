---
title: "Uses"
---

A running list of the hardware, software, and services I use day to day. Inspired by
[uses.tech](https://uses.tech).

## Editor & Terminal

Zed is my main editor these days — it's where I do real project work. Neovim is still always
open for quick edits and anything where I just need to get in and out fast. Used to lean on
Neovim for everything, but Zed won me over gradually.

Alacritty as my terminal, tmux for session management, fish as my shell.

## Languages & Tooling

Go for most of what I build. It's become my default when I'm starting something new.

React for frontends — TanStack was the most recent stack I tried before landing here. PHP and
Laravel power [TattooReserve](/projects/tattooreserve/), which is still running in production.

Terraform/OpenTofu for infrastructure, Ansible for configuration management. Docker or
Kubernetes depending on what the workload actually warrants — I try not to reach for
orchestration until there's a reason.

## Infrastructure

AWS for cloud workloads. Fly.io for quick projects and kicking tires on ideas — I've spent
enough time there to have opinions about it.

Ubuntu across all my VMs. Proxmox on two Dell PowerEdge servers running everything. Forgejo
for self-hosted Git, mirrored to GitHub as a backup. Woodpecker CI as my primary pipeline
runner with GitHub Actions kept in sync but normally disabled. n8n for lightweight workflows.
Jenkins for anything that needs a proper pipeline.

## Hardware

Custom Linux desktop running Arch Linux with i3wm on dual 1440p monitors.

Three Dell PowerEdge servers — one running Unraid, two running Proxmox. More homelab than I
planned to have, exactly as much as I need.
