package scanner

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"sync"

	"github.com/yanmengfei/poc-engine-soc/build"
)

type scanner struct {
	poc *build.PocEvent
}

var scannerPool = sync.Pool{New: func() interface{} { return new(scanner) }}

func (s *scanner) Start(target string) (verufy bool, err error) {
	cmd := exec.Command(s.poc.Python, "-c", s.poc.GetExecCode(target))
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	defer s.poc.CleanCache()
	if cmd.Run() != nil {
		return false, errors.New(stderr.String())
	}
	switch strings.TrimSpace(stdout.String()) {
	case "True":
		return true, nil
	case "False":
		return false, nil
	}
	return false, errors.New(stdout.String())
}

func New(poc *build.PocEvent) (scan *scanner, err error) {
	if poc == nil {
		return nil, errors.New("invalid poc")
	}
	scan = scannerPool.Get().(*scanner)
	scan.poc = poc
	return scan, nil
}

func Release(scan *scanner) {
	scan.poc = nil
	scannerPool.Put(scan)
}
