package clickhouse

// replicaStatus is the internal view over the one system.replicas row List
// reads, so the warnings stay a pure, testable function.
type replicaStatus struct {
	Readonly       uint8
	AbsoluteDelay  uint64
	ActiveReplicas uint32
	TotalReplicas  uint32
	ReadonlyCount  uint64
	// ReadonlyTables names the first few read-only tables, not all of them.
	ReadonlyTables string
	MaxDelay       uint64
	StuckCount     uint64
	StuckError     string
}
