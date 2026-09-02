package clients

import "time"

// IntegrationTimeout is the shared network timeout for clients that talk to
// same-host / local-network services (etcd, postgres, patroni's HTTP API). A
// healthy round-trip on that kind of deployment is well under 100ms, so this
// stays short enough to fail fast on an unreachable node instead of stalling
// requests for several seconds.
const IntegrationTimeout = 300 * time.Millisecond

// SshTimeout bounds an ssh connection attempt. It is its own constant because
// ssh is the one client that reaches a machine the operator chose rather than a
// service beside the database - a remote datacenter or a VPN hop away - and
// x/crypto/ssh spends it on the key exchange as well as the dial, so it has to
// cover several round trips, not one.
const SshTimeout = 5 * time.Second

// ExternalTimeout is the shared network timeout for clients that talk to
// external services reached over a real network rather than a local Docker
// deployment (LDAP, OIDC), where round-trips are naturally slower and a
// false timeout would fail a real user login.
const ExternalTimeout = 5 * time.Second

// CommandExecuteTimeout bounds a single one-shot remote command (docker
// run/exec/stop/start/list, ...) run through console.Command.Execute. It is
// generous enough to cover a slow image pull, but keeps a single hung
// command (e.g. a stuck "docker pull") from blocking its caller - and with
// it a whole synchronous cluster deploy (cluster.Service.Deploy) - forever.
// It does not apply to Start/Wait used directly for long-running log-follow
// streaming.
const CommandExecuteTimeout = 5 * time.Minute
