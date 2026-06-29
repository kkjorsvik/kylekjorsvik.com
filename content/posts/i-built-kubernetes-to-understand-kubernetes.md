---
title: "I Built Kubernetes to Understand Kubernetes"
date: "2026-06-28"
tags: ["go", "kubernetes", "devops", "networking"]
draft: false
description: "I'd set up Kubernetes clusters in my homelab and spent real time learning it, but still couldn't honestly explain what it was doing under the hood. So I built my own orchestrator — Smith — to earn the understanding by feeling every decision from the inside."
---

I'd set up Kubernetes more than once — clusters in my homelab, and a proof-of-concept at work
that got shelved the moment higher-priority work landed — and I'd put real time into learning
how to drive it. I still couldn't honestly tell you what it was doing under the hood.

I could write the YAML. I could `kubectl get pods`, read the events, fix the
CrashLoopBackOff, scale a Deployment, debug an Ingress that 404'd. I knew the nouns. But
if you'd stopped me and asked "what actually happens, kernel to kernel, when a pod on node
A talks to a Service backed by a pod on node B?" — I would've waved my hands at "kube-proxy
does some iptables stuff" and hoped you didn't ask a follow-up.

That gap bothered me more than I wanted to admit. I'd learned to drive the system without
ever understanding it — the internals were still magic. Reading docs and standing up clusters
taught me the operator's view; it didn't teach me the machine underneath. So I decided to
build my own — not to replace Kubernetes, not to fix anything about it, but to earn the
understanding by feeling every decision it makes from the inside. I called it Smith.

## Why not just read the source

The obvious objection: Kubernetes is open source. Go read it.

I tried. Reading the source teaches you what the code does. It does not teach you why every
other design was worse. You read kube-proxy's iptables mode and it looks arbitrary — a pile
of chains and a `--masquerade-all` flag that seems heavy-handed. Reading it, you nod along.
You don't understand it.

You understand it the day your own service load balancer silently drops every request where
the chosen backend happens to live on the same node as the caller, and you spend an evening
with tcpdump discovering why you suddenly need to masquerade traffic you thought you had no
reason to touch. Then you go back and reread that flag and realize it wasn't heavy-handed.
It was someone else's identical evening, encoded.

Reading is how you collect other people's conclusions. Building is how you earn your own. I
wanted the scar tissue.

## What Smith is

Smith is a multi-node container orchestrator written in Go. It's about 30 source files —
small enough to hold in your head, which is the entire point.

The shape will look familiar if you know Kubernetes:

- A control plane that holds desired state in SQLite and runs a reconcile loop.
- Agents on each node that own the local container lifecycle on top of containerd.
- mTLS between every node — TLS 1.3, `RequireAndVerifyClientCert`, all signed by a Smith CA.
  Agents authenticate to the control plane and the control plane authenticates to agents;
  there's no unauthenticated path on the internal plane.
- CNI networking — each node runs a `smith0` bridge with host-local IPAM. The control plane
  carves a `10.22.0.0/16` pool into a `/24` per node and installs static routes so containers
  route across nodes by their real IPs.
- iptables service load balancing — ClusterIPs and NodePorts, kube-proxy's iptables mode
  rebuilt small: a dispatch chain that jumps to a per-service chain, which DNATs across
  backends using the statistic module with probability `1/(N-i)` so the distribution stays
  uniform.
- Ingress with TLS termination — agents run a reverse proxy on `:443`, terminate with a
  wildcard cert (Let's Encrypt via the ACME DNS-01 challenge against Route 53), route by Host
  to the right ClusterIP, and redirect `:80`. Certs hot-swap on renewal without dropping the
  listener.
- Rolling updates — the reconciler detects spec drift by hashing the container-defining
  fields and replaces replicas within a MaxUnavailable budget.

Right now it runs on Proxmox: one control-plane VM plus five agent nodes, with real DNS
records in Route 53 pointing at it. It is genuinely serving traffic, not just passing tests.

![The Smith control-plane dashboard: five agent nodes alive, with deployops and postgres workloads running on real pod IPs, ClusterIPs, NodePorts, and a host-based ingress.](/static/images/smith-cluster-web-dashboard.png)
*The control-plane dashboard — five agents alive, serving `deployops` and `postgres` over a real ingress.*

![Proxmox showing the Smith cluster VMs: one control plane (Smith-Server-01) and five agents (Smith-Agent-01 through 05).](/static/images/proxmox-smith-cluster-nodes.png)
*The cluster on Proxmox: one control plane plus five agents.*

## The moments that actually taught me something

### 1. Containers on different nodes don't talk, and nobody tells you that

The first time I got CNI working, I was thrilled. A container came up, got `10.22.1.7`,
could reach the internet. Two containers on the same node could reach each other over the
bridge. I thought networking was done.

Then I scheduled two containers on different nodes and nothing worked, with no error. CNI
gives each node an isolated bridge. That's all it gives you. Nothing routes between those
bridges — that's explicitly out of scope for the plugin, and I'd never noticed because in
Kubernetes that's the CNI plugin's job, the part I'd always treated as a checkbox in the
install docs.

So I wrote the routing layer myself: a manager that installs a static route per peer node's
`/24` via the node's IP and reconciles them on a ticker, deleting routes for nodes that go
away. Then I hit the part that actually humbled me. I'd left `ipMasq` on in the bridge
config, the sane-looking default — so every container's traffic got masqueraded to the node
IP, and cross-node connections all appeared to originate from the node, not the pod.
Anti-affinity, source-based anything, debugging — all broken. The fix was to turn masquerade
off in CNI and do it selectively in the firewall, only where it's actually required. Which
led directly to:

### 2. The hairpin packet that ruined a weekend

A pod calls a Service. The load balancer's random pick lands on a backend that happens to be
a pod on the same node as the caller. Request goes out, never comes back. Every other case
worked. This one silently hung.

What's happening: the packet gets DNAT'd to the local backend and delivered straight over
the bridge. The backend replies directly to the caller — also local, also over the bridge —
so the reply never passes back through the node's conntrack, never gets un-DNAT'd, and
arrives at the caller as a packet from an address it never sent to. The caller drops it.
Correctly.

The fix is to mark service traffic and MASQUERADE it (fwmark `0x4000`) so the reply is
forced back through conntrack and rewritten on the way out. The cost is that backends see
the node IP instead of the client pod IP. That is exactly what `--masquerade-all` does in
kube-proxy. I had read that flag a dozen times. I only understood it after building the bug
it prevents.

### 3. "I scheduled it" is not "it's running"

My first reconcile loop was edge-triggered in spirit: figure out where a workload goes, POST
it to the agent, record that I'd done it, move on. Clean. Also wrong.

Nodes reboot. Agents restart. A push succeeds and the container dies thirty seconds later.
The control plane's memory said "placed" while reality said "gone," and nothing reconverged.
The whole thing only works once you accept that push once is a lie and the loop has to
continuously compare desired against the cluster's actual running state — fetched fresh from
each agent every tick — and re-push anything that isn't where it should be. With a grace
period, so you don't murder a container that's simply still starting.

That's the lesson the word "reconcile" is hiding. It's not "apply changes." It's "assume
everything you believe is already stale and prove it again, forever." Level-triggered, not
edge-triggered. I'd read that phrase a hundred times and it meant nothing until my own naive
version rotted in front of me.

### 4. Rolling updates with no one in charge

Replacing a replica makes it briefly unavailable. Replace too many at once and you take an
outage. Kubernetes makes this look trivial. Doing it without a central lock is not.

My reconciler recomputes, on every single tick, how many replicas of a workload are
currently not running, and only rolls a replica to a new spec if doing so keeps that count
under the workload's MaxUnavailable. No coordinator, no lease, no lock — just the same
level-triggered comparison applied to availability, where a partially-completed rollout is a
perfectly valid state to be discovered and continued next tick. Getting comfortable with
"interrupted halfway is fine, I'll figure it out from observed state when I wake up" was the
real unlock, and it's the same idea as #3 wearing different clothes.

## What Kubernetes got right that I didn't appreciate before

Building a worse version of something is the most honest way to respect the original. A few
things I now believe in for reasons instead of by reputation:

- Level-triggered reconciliation is the whole game. It's not a style choice. It is the only
  model that survives contact with nodes that reboot and networks that partition. Everything
  robust about Kubernetes falls out of this one decision, and I had to build the broken
  edge-triggered version first to see it.
- The CNI boundary is a gift. The line between "give a container an IP on a bridge" and
  "make bridges reach each other" is drawn exactly where it should be. I resented it until I
  understood it was protecting the plugin from my routing problems.
- etcd earns its complexity. Smith keeps its brain in a single SQLite file on a single
  control-plane node. That's a real single point of failure, and I'm running it knowingly. A
  replicated consensus store is not enterprise bloat — it's the answer to a problem I now
  have and am choosing to defer.
- Probes gate rollouts for a reason. Smith health-checks workloads (HTTP and exec, with
  failure thresholds) and feeds that back into the loop. "Is the process up" and "is it ready
  to serve" are different questions, and conflating them is how you roll a deploy straight
  into an outage.
- `--masquerade-all` was never laziness. See above. I owe that flag an apology.

The humbling summary: almost everything in Kubernetes I'd privately filed under
"overengineered" turned out to be load-bearing. The parts that look like cruft are usually
someone's scar tissue, same as mine, just older.

## Where Smith is going

Smith already does the unglamorous operational stuff that makes me trust it: subnet
allocations persist in SQLite and are keyed by node ID, so a rebooting node gets the same
`/24` back instead of having it reassigned out from under its running containers. Desired
state lives in a git repo and applies declaratively through a small CLI, with secrets
encrypted at rest. It's boring on purpose.

![smithctl apply --dry-run against the GitOps repo, showing the planned workloads, services, and ingress, with secret env vars shown only as "set from overlay".](/static/images/smith-apply-dry-run-gitops.png)
*`smithctl apply --dry-run` against the GitOps repo — secrets come from an encrypted overlay, surfaced only as "set from overlay".*

Next is getting it onto real hardware and real load. The control plane and agents are moving
off Proxmox VMs onto two Dell PowerEdge R815 nodes, with workers spread across both, to run
actual homelab workloads instead of test pods. That's the real test — not "does it pass CI"
but "do I trust it with things I'd be annoyed to lose." (Git and CI/CD deliberately stay
outside Smith; the orchestrator's job is running workloads, not being my build plane.)

## Closing

I didn't build Smith to ship a product. I built it because it gnawed at me to have
spent so much time learning to operate a tool I still couldn't explain from the inside, and
the only cure I've ever found for that is to build the thing badly myself until the good
version's decisions become obvious.

It worked. I can now answer the node-A-to-node-B question in detail, with the specific
reasons each layer exists, because I broke each one personally. That's a different kind of
knowledge than I had before, and I want more jobs that demand it.

So, plainly: I'm open to Senior DevOps and Platform Engineering roles. If your team runs
infrastructure where understanding why the system behaves the way it does is the actual job
— not just keeping it green — that's the work I like and the way I like to learn it.
[Smith is on GitHub](https://github.com/kkjorsvik/smith)
