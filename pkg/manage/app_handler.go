package manage

import (
	"context"

	"github.com/go-faster/sdk/zctx"
	"go.uber.org/zap"
)

// Compile-time check for Handler.
var _ Handler = (*Entrypoint)(nil)

type Entrypoint struct {
	UnimplementedHandler // automatically implement all methods
}

func (h Entrypoint) GetRaftNode(ctx context.Context, params GetRaftNodeParams) (*GetRaftNodeOK, error) {
	log := zctx.From(ctx)

	log.Info("GetRaftNode", zap.Any("params", params))
	return &GetRaftNodeOK{
		Nodes: []RaftMetaNode{
			{
				ID:     "node1",
				Addr:   "1",
				Status: RaftStatusFollower,
				Role:   RaftRoleVoter,
			},
		},
	}, nil
}
