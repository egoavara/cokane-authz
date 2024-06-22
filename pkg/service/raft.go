package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
	openfga "github.com/openfga/go-sdk"
)

type RaftService struct {
	Raft *raft.Raft
}

func NewRaftService(r *raft.Raft) *RaftService {
	return &RaftService{
		Raft: r,
	}
}

func (s *RaftService) Middleware(engine *gin.Engine) {
	router := engine.Group("/raft")
	router.GET("/join", s.join)
	router.GET("/snapshot", func(ctx *gin.Context) {
		fut := s.Raft.Snapshot()
		if err := fut.Error(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.Status(http.StatusNoContent)
	})
	router.GET("/setup_model", func(ctx *gin.Context) {
		cmd, err := json.Marshal(&OpenFgaFSMCommand{
			Type: OpenFgaFSMCommandTypeModel,
			Model: &openfga.AuthorizationModel{

				SchemaVersion: "1.1",
				TypeDefinitions: []openfga.TypeDefinition{
					openfga.TypeDefinition{
						Type: "user",
					},
				},
			},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fut := s.Raft.Apply(cmd, 2*time.Second)
		if err := fut.Error(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp := fut.Response()
		fmt.Println("Response:", resp)
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

}

type joinQuery struct {
	followerId   string `form:"followerId" validate:"required, min=1, max=100, alphanum"`
	followerAddr string `form:"followerAddr" validate:"required, hostname_port"`
}

func (s *RaftService) join(ctx *gin.Context) {
	var query joinQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if s.Raft.State() != raft.Leader {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Not a leader"})
		return
	}

	err := s.Raft.AddVoter(raft.ServerID(query.followerId), raft.ServerAddress(query.followerAddr), 0, 0).Error()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
