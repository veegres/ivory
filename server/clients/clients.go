package clients

import "time"

// IntegrationTimeout is the shared network timeout for clients that talk to
// same-host / local-network services (etcd, postgres, ssh, patroni's HTTP
// API). A healthy round-trip on that kind of deployment is well under
// 100ms, so this stays short enough to fail fast on an unreachable node
// instead of stalling requests for several seconds.
const IntegrationTimeout = 300 * time.Millisecond

// ExternalTimeout is the shared network timeout for clients that talk to
// external services reached over a real network rather than a local Docker
// deployment (LDAP, OIDC), where round-trips are naturally slower and a
// false timeout would fail a real user login.
const ExternalTimeout = 5 * time.Second
