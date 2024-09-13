// Code generated by ogen, DO NOT EDIT.

package manage

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// GetRaftNode implements getRaftNode operation.
//
// Get raft node information.
//
// GET /raft/node
func (UnimplementedHandler) GetRaftNode(ctx context.Context, params GetRaftNodeParams) (r *GetRaftNodeOK, _ error) {
	return r, ht.ErrNotImplemented
}

// JoinRaftCluster implements joinRaftCluster operation.
//
// Join a raft cluster.
//
// POST /raft/node
func (UnimplementedHandler) JoinRaftCluster(ctx context.Context, req *RaftBaseNode) (r *RaftMetaNode, _ error) {
	return r, ht.ErrNotImplemented
}