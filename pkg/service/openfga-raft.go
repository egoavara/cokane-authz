package service

import (
	"fmt"

	"egoavara.net/authz/pkg/config"
	"github.com/go-resty/resty/v2"
)

type (
	OpenfgaRaftService struct {
		BootstrapAddrs []string
		Addrs          []string
	}
	OpenfgaRaftJoinRequest struct {
		ID   string `json:"id"`
		Addr string `json:"addr"`
	}
	OpenfgaRaftNodeList struct {
		Nodes []OpenfgaRaftNode `json:"nodes"`
	}
	OpenfgaRaftNode struct {
		ID   string `json:"id"`
		Addr string `json:"addr"`
	}
)

func NewOpenfgaRaftService(bootstrapAddrs []string) *OpenfgaRaftService {
	return &OpenfgaRaftService{
		BootstrapAddrs: bootstrapAddrs,
		Addrs:          nil,
	}
}

func (s *OpenfgaRaftService) Join(request *OpenfgaRaftJoinRequest) error {
	resp, err := resty.New().R().SetBody(request).Post(fmt.Sprintf("http://auth.egoavara.net/raft/%s/join", config.Metadata.Value().Version))

	return nil
}

func (s *OpenfgaRaftService) Leader() (*OpenfgaRaftNode, error) {
	resty.New().R().Get("http://" + s.Addrs[0] + "/raft/${}/leader")
	return nil, nil
}

func (s *OpenfgaRaftService) List() (*OpenfgaRaftNodeList, error) {
	return s.list()
}
