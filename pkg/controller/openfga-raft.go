package controller

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"egoavara.net/authz/pkg/fsm"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	openfga "github.com/openfga/go-sdk"
)

type (
	OpenfgaRaftConfig struct {
		Cluster OpenfgaRaftClusterConfig `json:"raft" yaml:"raft" validate:"required"`
		Openfga OpenfgaRaftOpenfgaConfig `json:"openfga" yaml:"openfga" validate:"required"`
	}
	OpenfgaRaftClusterConfig struct {
		BootstrapAddrs  []string `json:"bootstrapAddrs" yaml:"bootstrapAddrs"`
		LogStoreFile    string   `json:"logStoreFile" yaml:"logStoreFile" validate:"required"`
		StableStoreFile string   `json:"stableStoreFile" yaml:"stableStoreFile" validate:"required"`
		SnapshotDir     string   `json:"snapshotDir" yaml:"snapshotDir" validate:"required"`
	}
	OpenfgaRaftOpenfgaConfig struct {
		Addr string `json:"addr" yaml:"addr" validate:"required"`
	}
)

type OpenfgaRaft struct {
	config     *OpenfgaRaftConfig
	FgaClient  *openfga.APIClient
	FSM        *fsm.OpenFgaFSM
	Transfort  *raft.NetworkTransport
	Raft       *raft.Raft
	RaftConfig *raft.Config
}

func NewOpenfgaRaft(config *OpenfgaRaftConfig) (*OpenfgaRaft, error) {
	raft := new(OpenfgaRaft)
	if err := raft.Reload(config); err != nil {
		return nil, err
	}
	return raft, nil
}

func (s *OpenfgaRaft) Reload(config *OpenfgaRaftConfig) error {
	// OpenFga 서비스 설정
	fgaConfig, err := openfga.NewConfiguration(openfga.Configuration{
		ApiUrl: config.Openfga.Addr,
		RetryParams: &openfga.RetryParams{
			MaxRetry: 3,
		},
	})
	if err != nil {
		return err
	}
	fgaClient := openfga.NewAPIClient(fgaConfig)
	// Raft 서비스 설정
	selfHostname := s.selfHostname()
	selfPort := 9812
	selfAddr := fmt.Sprintf("%s:%d", selfHostname, selfPort)
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(selfHostname)
	logStore, err := raftboltdb.NewBoltStore(config.Cluster.LogStoreFile)
	if err != nil {
		return err
	}
	stableStore, err := raftboltdb.NewBoltStore(config.Cluster.StableStoreFile)
	if err != nil {
		return err
	}
	snapshotStore, err := raft.NewFileSnapshotStore(config.Cluster.SnapshotDir, 1, os.Stderr)
	if err != nil {
		return err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", selfAddr)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(selfAddr, tcpAddr, 10, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}
	fsm := fsm.NewOpenFgaFSM(fgaClient)
	raftNode, err := raft.NewRaft(raftConfig, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return err
	}
	shouldBootstrap, err := s.shouldBootstrap(config, selfAddr, logStore, stableStore, snapshotStore)
	if err != nil {
		return err
	}
	if shouldBootstrap {
		configuration := raft.Configuration{}
		for _, addr := range config.Cluster.BootstrapAddrs {
			configuration.Servers = append(configuration.Servers, raft.Server{
				ID:      raft.ServerID(addr),
				Address: raft.ServerAddress(addr),
			})
		}
		raftNode.BootstrapCluster(configuration)
	} else {

	}
	// 완료된 설정 데이터 저장
	s.config = config
	s.FgaClient = fgaClient
	s.FSM = fsm
	s.Transfort = transport
	s.Raft = raftNode
	s.RaftConfig = raftConfig
	return nil
}
func (s *OpenfgaRaft) shouldBootstrap(config *OpenfgaRaftConfig, selfAddr string, logStore raft.LogStore, stableStore raft.StableStore, snapsStore raft.SnapshotStore) (bool, error) {
	exists, err := raft.HasExistingState(logStore, stableStore, snapsStore)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	if len(config.Cluster.BootstrapAddrs) == 0 {
		// 부트스트랩 주소가 없는 경우 (일반적으로 잘못된 설정으로 인해 발생)
		return false, nil
	}
	first := config.Cluster.BootstrapAddrs[0]
	firstTCPAddr, err := net.ResolveTCPAddr("tcp", first)
	if err != nil {
		return false, err
	}
	selfTCPAddr, err := net.ResolveTCPAddr("tcp", selfAddr)
	if err != nil {
		return false, err
	}
	if !(firstTCPAddr.IP.Equal(selfTCPAddr.IP) && firstTCPAddr.Port == selfTCPAddr.Port && firstTCPAddr.Zone == selfTCPAddr.Zone) {
		return false, nil
	}
	return true, nil
}

func (s *OpenfgaRaft) selfHostname() string {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	return hostname
}

func (s *OpenfgaRaft) UseRouter(engine *gin.Engine, router gin.IRouter) {
	router.GET("/leader", s.leader)
	router.POST("/join", s.join)
	router.POST("/snapshot", func(ctx *gin.Context) {
		fut := s.Raft.Snapshot()
		if err := fut.Error(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.Status(http.StatusNoContent)
	})
}

type JoinRole string

const (
	JoinRoleVoter    JoinRole = "voter"
	JoinRoleNonVoter JoinRole = "non-voter"
)

type (
	JoinQuery struct {
		Id   string   `form:"id" validate:"required, min=1, max=100, alphanum"`
		Addr string   `form:"addr" validate:"required, hostname_port"`
		Role JoinRole `form:"role" validate:"required, oneof=voter non-voter"`
	}
	NodeResponse struct {
		Id   raft.ServerID      `json:"id"`
		Addr raft.ServerAddress `json:"addr"`
	}
	NodeListResponse struct {
		Nodes []NodeResponse `json:"nodes"`
	}
	LeaderResponse NodeResponse
)

func (s *OpenfgaRaft) join(ctx *gin.Context) {
	var query JoinQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if s.Raft.State() != raft.Leader {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Not a leader"})
		return
	}
	if query.Role == JoinRoleNonVoter {
		err := s.Raft.AddNonvoter(raft.ServerID(query.Id), raft.ServerAddress(query.Addr), 0, 0).Error()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		err := s.Raft.AddVoter(raft.ServerID(query.Id), raft.ServerAddress(query.Addr), 0, 0).Error()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.Status(http.StatusNoContent)
}

func (s *OpenfgaRaft) leader(ctx *gin.Context) {
	addr, id := s.Raft.LeaderWithID()
	ctx.JSON(http.StatusOK, LeaderResponse{
		Id:   id,
		Addr: addr,
	})
}
