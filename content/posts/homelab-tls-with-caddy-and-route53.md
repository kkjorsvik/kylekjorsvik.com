---
title: "Homelab TLS with Caddy and Route 53"
date: "2026-04-12"
tags: ["caddy", "tls", "aws", "homelab"]
draft: false
description: "Wildcard certificates for internal-only services using Caddy's DNS challenge against Route 53 — no ports exposed to the internet."
---

I run a handful of services in my homelab that should never be reachable from the public
internet, but I still want real TLS certificates for them. Self-signed certs and the
endless browser warnings get old fast. The answer is the ACME DNS-01 challenge.

## The DNS-01 challenge

Unlike HTTP-01, the DNS-01 challenge proves domain control by creating a TXT record
instead of serving a file. That means you never have to expose port 80 or 443 — perfect
for internal services.

Caddy with the Route 53 DNS plugin handles the whole dance:

```caddyfile
*.home.example.com {
	tls {
		dns route53
	}
	reverse_proxy grafana.internal:3000
}
```

## Scoping the IAM permissions

Give Caddy an IAM user with only the Route 53 permissions it needs for the hosted zone:

```json
{
  "Effect": "Allow",
  "Action": [
    "route53:ListHostedZonesByName",
    "route53:GetChange",
    "route53:ChangeResourceRecordSets"
  ],
  "Resource": "*"
}
```

The result: a wildcard cert that renews itself, served to services that are only resolvable
on my LAN. Zero inbound firewall rules.
