package scanner

import (
	"context"

	"github.com/yanmengfei/poc-engine-soc/build"
	"github.com/yanmengfei/poc-engine-soc/proto"
)

type Scanner struct {
	poc *build.SocpocEvent
	// client proto.SocpocClient
}

func New(poc *build.SocpocEvent) *Scanner {
	return &Scanner{poc: poc}
}

func (s *Scanner) Start(target string, ctx context.Context) (bool, error) {
	client := proto.NewSocpocClient(s.poc.Conn)
	resp, err := client.Execute(ctx, &proto.ExecuteRequest{Key: s.poc.Key, Module: s.poc.Module, Url: target})
	if err != nil {
		return false, err
	}
	return resp.Status, err
}
