package etcd

// member and endpointStatus are internal views over clientv3 types so the
// response mapping stays a pure, testable function.

type member struct {
	ID         uint64
	Name       string
	ClientURLs []string
	IsLearner  bool
}

type endpointStatus struct {
	Leader    uint64
	RaftIndex uint64
	Err       error
}

type switchoverBody struct {
	Leader      *string `json:"leader"`
	Candidate   *string `json:"candidate"`
	ScheduledAt *string `json:"scheduled_at"`
}
