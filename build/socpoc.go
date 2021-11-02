package build

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"
	"time"

	"github.com/yanmengfei/poc-engine-soc/proto"
	"google.golang.org/grpc"
)

type SocpocEvent struct {
	Conn   *grpc.ClientConn
	Key    string
	Module string
}

const (
	startIndexStr = "class "
	stopIndexStr  = "(BasicPoc):"
)

func (e *SocpocEvent) dial(addr string) (err error) {
	e.Conn, err = grpc.Dial(
		addr, grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
		grpc.FailOnNonTempDialError(true),
	)
	return err
}

func (e *SocpocEvent) make(code string) error {
	hash := md5.New()
	scan := bufio.NewScanner(strings.NewReader(code))
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, startIndexStr) && strings.HasSuffix(line, stopIndexStr) {
			e.Module = line[strings.Index(line, " ")+1 : strings.Index(line, stopIndexStr)]
		}
		io.WriteString(hash, line)
	}
	e.Key = hex.EncodeToString(hash.Sum(nil))
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	client := proto.NewSocpocClient(e.Conn)
	defer cancel()
	_, err := client.Setenv(ctx, &proto.SetenvRequest{Key: e.Key, Code: code})
	return err
}

func NewPocEvent(addr string, code string) (*SocpocEvent, error) {
	event := SocpocEvent{}
	if err := event.dial(addr); err != nil {
		return nil, err
	}
	if err := event.make(code); err != nil {
		return nil, err
	}
	return &event, nil
}
