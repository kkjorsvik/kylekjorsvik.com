---
title: "Docker Fleet Management Without Kubernetes"
date: "2026-03-03"
tags: ["docker", "devops", "ansible"]
draft: false
description: "Not every workload needs Kubernetes. How I manage a small fleet of Docker hosts with Compose, Ansible, and a bit of discipline."
---

Kubernetes is a remarkable piece of engineering, and for a lot of teams it's also wildly
overkill. For a handful of single-tenant services on a few VMs, the operational tax of a
cluster — etcd, upgrades, CNI, the YAML — buys you very little.

## The setup

My approach for small fleets:

- **Docker Compose** as the unit of deployment per host
- **Ansible** to template Compose files and push them out
- **Watchtower** (carefully scoped) for image updates on non-critical services
- **A shared overlay network** only where services actually need to talk

The deploy is a single playbook run:

```yaml
- name: Deploy stack
  hosts: app_servers
  tasks:
    - name: Template compose file
      template:
        src: docker-compose.yml.j2
        dest: /opt/stack/docker-compose.yml

    - name: Bring up the stack
      community.docker.docker_compose_v2:
        project_src: /opt/stack
        pull: always
```

## When to graduate

You graduate to Kubernetes (or Nomad) when you need real bin-packing, autoscaling, or
multi-tenant isolation. Until then, boring infrastructure is a feature.
