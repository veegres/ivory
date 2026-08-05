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
	// Self is true only on the entry describing the connection's own node,
	// used to find "am I currently the primary" without a separate command.
	Self bool `bson:"self"`
}
