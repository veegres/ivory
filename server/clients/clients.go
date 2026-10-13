package clients

import "time"

// IntegrationDialTimeout bounds only the tcp connect to a database or keeper
// node, so an unreachable node still fails fast. It applies where Ivory owns
// the dialer; a driver whose dial also covers tls and auth gets
// IntegrationTimeout instead, since those take several round trips.
const IntegrationDialTimeout = 1 * time.Second

// IntegrationTimeout bounds a whole request to a database or keeper node. It is
// not a local-network budget: Ivory may sit a VPN hop away from the cluster, and
// a request it abandons keeps running on the node, so giving up too early piles
// work onto a keeper rather than sparing it.
const IntegrationTimeout = 5 * time.Second

// SshTimeout bounds an ssh connection attempt. x/crypto/ssh spends it on the
// key exchange as well as the dial, so it has to cover several round trips.
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
