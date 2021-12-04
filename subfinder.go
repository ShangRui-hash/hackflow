package hackflow

import (
	"bufio"
	"fmt"
	"go/build"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Subfinder struct {
	BaseTool
}

func newSubfinder() Tool {
	return &Subfinder{
		BaseTool: BaseTool{
			name: SUBFINDER,
			desp: "被动子域名收集工具",
		},
	}
}

func GetSubfinder() *Subfinder {
	return container.Get(SUBFINDER).(*Subfinder)
}

func (s *Subfinder) ExecPath() (string, error) {
	return s.BaseTool.ExecPath(s.Download)
}

func (s *Subfinder) Download() (string, error) {
	if err := GetGo().Install("github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"); err != nil {
		return "", err
	}
	return build.Default.GOPATH + "/bin/subfinder", nil
}

type SubfinderRunConfig struct {
	DomainCh                       chan string
	Proxy                          string `flag:"-proxy"`
	Domain                         string `flag:"-d"`
	RoutineCount                   int    `flag:"-t"`
	RemoveWildcardAndDeadSubdomain bool   `flag:"-nW"`
	OutputInHostIPFormat           bool   `flag:"-oI"`
	OutputInJsonLineFormat         bool   `flag:"-oJ"`
}

func (s *Subfinder) Run(config *SubfinderRunConfig) (subdomainCh chan string, err error) {
	args := append([]string{"-silent", "-nC"}, parseConfig(*config)...)
	execPath, err := s.ExecPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Error("cmd.StdoutPipe failed,err:", err)
		return nil, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logrus.Error("cmd.StdinPipe failed,err:", err)
		return nil, err
	}
	if config.DomainCh != nil {
		go func() {
			for domain := range config.DomainCh {
				fmt.Fprintln(stdin, domain)
			}
			stdin.Close()
		}()
	}
	if err := cmd.Start(); err != nil {
		logrus.Error("Execute failed when Start:", err)
		return nil, err
	}
	logger.Debugf("%s 启动成功\n", s.name)
	subdomainCh = make(chan string, 1024)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			subdomainCh <- scanner.Text()
		}
		close(subdomainCh)
	}()
	return subdomainCh, nil
}
