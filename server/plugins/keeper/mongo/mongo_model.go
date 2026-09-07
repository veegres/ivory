package mongo

import "time"

// replSetStatus mirrors the subset of the replSetGetStatus command's reply
// this adapter reads. Any member (primary or secondary) can report the full
// set, the same way Patroni's /cluster or etcd's member list do, since it is
// built from the target's own heartbeat view of every peer.
type replSetStatus struct {
	Set     string          `bson:"set"`
	Members []replSetMember `bson:"members"`
}

type replSetMember struct {
	Name       string    `bson:"name"`
	StateStr   string    `bson:"stateStr"`
	Health     float64   `bson:"health"`
	OptimeDate time.Time `bson:"optimeDate"`
	// SyncSourceHost is the member this one actually replicates from, which is
	// not necessarily the primary - mongo lets a secondary chain off another
	// secondary, so a whole branch of the set can fall behind together.
	SyncSourceHost string `bson:"syncSourceHost"`
	// PingMs is the round trip time of the last heartbeat, reported only for
	// peers - the member describing itself has nothing to ping.
	PingMs int64 `bson:"pingMs"`
	// Self is true only on the entry describing the connection's own node,
	// used to find "am I currently the primary" without a separate command.
	Self bool `bson:"self"`
}
