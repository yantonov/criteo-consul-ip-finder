# AGENTS.md

## Project
Consul IP Finder — CLI tool that finds Consul service names by IP address across datacenters and environments.
Stack: Go 1.23, standard library, no database.

## Starting a session
Work lives in `bl` tickets. Run `bl help` for the command surface. If `.tickets/` doesn't exist yet, bootstrap with `bl init`.

## Definition of Done
A feature is done when `make check` is green, the feature's own verification passes, and the ticket is closed with `bl done <id>`.

## Hard constraints
1. MUST NOT commit secrets; all values come from CLI flags or environment variables.
2. MUST NOT change the Consul relay URL pattern (`consul-relay.{dc}.{env}.crto.in`) without explicit human approval — it is an infrastructure contract.

## Work rules
- **WIP=1.** Claim one ticket at a time: `bl start <id>`, close it with `bl done <id>` before claiming another.
- Stay inside the claimed ticket. Record any other problem with `bl add "<title>"` and continue.
- Priority order: correctness → performance → style.

## Escalate to human (stop and ask)
- Changing the Consul relay URL pattern.
- New external dependency.
- Any decision affecting security or privacy.
- Deleting anything that cannot be recovered.
