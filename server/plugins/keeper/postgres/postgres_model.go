package postgres

// nodeStats is the internal view over the columns listQuery reads beside the
// node's role and lag, so the response mapping stays a pure, testable function.
type nodeStats struct {
	Version        string
	Connections    int
	MaxConnections int
}
