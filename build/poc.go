package build

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yanmengfei/poc-engine-soc/library/utils"
)

type PocEvent struct {
	Name   string
	Module string
	Code   string
	Cache  string
	Python string
}

const (
	runpath       = "pocs/"
	startIndexStr = "class "
	stopIndexStr  = "(BasicPoc):"
	template      = `import importlib;poc=getattr(importlib.import_module("pocs.%s"),"%s")();print(poc.run("%s"))`
)

func init() {
	if !utils.Exists(runpath) {
		os.Mkdir(runpath, os.ModePerm)
	}
}

func (p *PocEvent) GetExecCode(url string) string {
	return fmt.Sprintf(template, p.Name, p.Module, url)
}

func (p *PocEvent) CleanCache() error {
	return os.Remove(p.Cache)
}

func NewPocEvent(id, code, python string) (*PocEvent, error) {
	var module string
	scan := bufio.NewScanner(strings.NewReader(code))
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, startIndexStr) && strings.HasSuffix(line, stopIndexStr) {
			module = line[strings.Index(line, " ")+1 : strings.Index(line, stopIndexStr)]
			break
		}
	}
	if module == "" {
		return nil, errors.New("invalid poc")
	}
	var p = PocEvent{Cache: runpath + id + ".py"}
	if err := ioutil.WriteFile(p.Cache, utils.StrToBytes(code), fs.ModePerm); err != nil {
		return nil, err
	}
	p.Name = id
	p.Code = code
	p.Module = module
	p.Python = python
	return &p, nil
}
