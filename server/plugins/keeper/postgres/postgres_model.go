package postgres

// nodeStats is the internal view over the columns listQuery reads beside the
// node's role and lag, so the response mapping stays a pure, testable function.
type nodeStats struct {
	Version        string
	Connections    int
	MaxConnections int
	// Receiving is whether a standby's WAL receiver is running at all. Its row is
	// visible to every role, unlike the receiver's details.
	Receiving bool
	// UnreservedSlots and LostSlots name the slots by postgres' own wal_status.
	UnreservedSlots []string
	LostSlots       []string
}
