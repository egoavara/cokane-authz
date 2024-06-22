package acts

import (
	"github.com/hashicorp/raft"
	openfga "github.com/openfga/go-sdk"
)

type Act struct {
	OpenFga       *openfga.APIClient
	Raft          *raft.Raft
	LogStore      raft.LogStore
	StableStore   raft.StableStore
	SnapshotStore raft.SnapshotStore
	Transport     raft.Transport
}
