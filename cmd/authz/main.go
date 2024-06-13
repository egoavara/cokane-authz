package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"egoavara.net/authz/acts"
	"egoavara.net/authz/service"
	"egoavara.net/authz/util"
	"egoavara.net/authz/works"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	openfga "github.com/openfga/go-sdk"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	ctx := context.Background()
	hostname := util.MustHostname()

	// # OpenFga 사용
	// OpenFga 서비스에 연결합니다.
	fgaConfig := util.Must(openfga.NewConfiguration(openfga.Configuration{
		ApiUrl: "http://localhost:8080",
		RetryParams: &openfga.RetryParams{
			MaxRetry: 3,
		},
	}))
	fgaClient := openfga.NewAPIClient(fgaConfig)

	// # Raft 노드 관리
	// 클러스터 상태를 유지하기 위해 Raft 서비스를 시작합니다.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(hostname)
	logStore := util.Must(raftboltdb.NewBoltStore("./raft-log.db"))
	stableStore := util.Must(raftboltdb.NewBoltStore("./raft-stable.db"))
	snapshotStore := util.Must(raft.NewFileSnapshotStore("./raft-snapshot", 1, os.Stderr))
	raftAddress := fmt.Sprintf("%s:9812", hostname)
	tcpAddr := util.Must(net.ResolveTCPAddr("tcp", raftAddress))
	fmt.Println("Raft node started at:", tcpAddr.IP.IsUnspecified())
	transport := util.Must(raft.NewTCPTransport(raftAddress, tcpAddr, 10, 10*time.Second, os.Stderr))
	fsm := service.NewOpenFgaFSM(fgaClient)
	raftNode := util.Must(raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport))
	if true {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raft.ServerID(hostname),
					Address: transport.LocalAddr(),
				},
			},
		}
		raftNode.BootstrapCluster(configuration)
	}
	// # Temporal 사용
	// Temporal 서비스에 연결합니다.
	tClient := util.Must(client.Dial(client.Options{HostPort: client.DefaultHostPort}))
	defer tClient.Close()
	// Temporal 워커를 시작합니다.
	tWorker := worker.New(tClient, "query", worker.Options{
		BackgroundActivityContext: ctx,
	})
	tWorker.RegisterWorkflow(works.UserAcl)
	tWorker.RegisterActivity(&acts.Act{
		OpenFga:       fgaClient,
		Raft:          raftNode,
		LogStore:      logStore,
		StableStore:   stableStore,
		SnapshotStore: snapshotStore,
		Transport:     transport,
	})

	// # HTTP server 등록
	raftService := service.NewRaftService(raftNode)
	engine := gin.Default(
		raftService.Middleware,
	)
	engine.Run(":81")
}
