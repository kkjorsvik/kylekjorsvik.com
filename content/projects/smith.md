---
title: "Smith"
description: "A multi-node container orchestrator written in Go — built from scratch to understand Kubernetes from the inside."
status: "active"
tags: ["go", "kubernetes", "devops", "networking"]
draft: false
featured: true
links:
  - label: "GitHub"
    url: "https://github.com/kkjorsvik/smith"
  - label: "Blog post"
    url: "/blog/i-built-kubernetes-to-understand-kubernetes/"
---

Smith is a multi-node container orchestrator I wrote in Go to earn a real understanding of
how Kubernetes works — not by reading the source, but by building a smaller version of it and
feeling every design decision from the inside. It's about 30 source files, small enough to
hold in your head, and it genuinely serves traffic rather than just passing tests.

The shape will look familiar if you know Kubernetes: a control plane that holds desired state
in SQLite and runs a level-triggered reconcile loop, agents on each node that own the local
container lifecycle on top of containerd, and mTLS between every node signed by a Smith CA.

## What it does

- **Reconciliation** — a level-triggered control loop that continuously compares desired state
  against the cluster's actual running state and re-converges, surviving node reboots and
  partitions.
- **CNI networking** — a `smith0` bridge per node with host-local IPAM, a `10.22.0.0/16` pool
  carved into a `/24` per node, and static routes so containers reach each other across nodes
  by their real IPs.
- **Service load balancing** — ClusterIPs and NodePorts implemented with iptables, including
  selective masquerade to handle the hairpin case kube-proxy's `--masquerade-all` exists for.
- **Ingress with TLS** — agents run a reverse proxy on `:443`, terminating with a Let's Encrypt
  wildcard cert (ACME DNS-01 via Route 53) and hot-swapping certs on renewal.
- **Rolling updates** — spec drift detected by hashing container-defining fields, with replicas
  replaced inside a MaxUnavailable budget and no central lock.
- **GitOps** — desired state lives in a git repo and applies declaratively through a small CLI,
  with secrets encrypted at rest.

## Tech stack

Go, containerd, SQLite, iptables, CNI, mTLS (TLS 1.3), and ACME/Let's Encrypt. Currently runs
on Proxmox — one control-plane VM plus five agent nodes — with real DNS in Route 53.

## Status

Active. Smith is serving real traffic and moving off Proxmox VMs onto bare-metal Dell PowerEdge
hardware to run actual homelab workloads. The full story of what building it taught me is in the
[blog post](/blog/i-built-kubernetes-to-understand-kubernetes/).
