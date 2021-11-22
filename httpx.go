package hackflow

import (
	"bufio"
	"fmt"
	"go/build"
	"os/exec"
)

type Httpx struct {
	name     string
	execPath string
}

//Name 获取工具名
func (s *Httpx) Name() string {
	return s.name
}

//ExecPath 获取工具执行路径
func (s *Httpx) ExecPath() (string, error) {
	if s.execPath == "" {
		if err := s.download(); err != nil {
			logger.Errorf("download %s failed,err:%v", s.Name(), err)
			return "", err
		}
	}
	return s.execPath, nil
}

func (s *Httpx) download() error {
	if err := GetGo().Install("github.com/projectdiscovery/httpx/cmd/httpx@latest"); err != nil {
		return err
	}
	s.execPath = build.Default.GOPATH + "/bin/httpx"
	return nil
}

func GetHttpx() *Httpx {
	if tool := container.Get(HTTPX); tool != nil {
		return tool.(*Httpx)
	}
	container.Set(&Httpx{
		name: HTTPX,
	})
	return container.Get(HTTPX).(*Httpx)
}

type HttpxRunConfig struct {
	DisplaySatusCode     bool   `flag:"-sc"`
	DisplayContentLength bool   `flag:"-cl"`
	DisplayResponseTime  bool   `flag:"-rt"`
	DisplayTitle         bool   `flag:"-title"`
	DisplayRequestMethod bool   `flag:"-method"`
	DisplayHostIP        bool   `flag:"-ip"`
	DisplayHostName      bool   `flag:"-cname"`
	StoreHTTPResponse    bool   `flag:"-sr"`
	RandomAgent          bool   `flag:"-random-agent"`
	Threads              int    `flag:"-t"`
	RateLimit            int    `flag:"-rl"`
	Proxy                string `flag:"-proxy"`
	URLCh                chan string
}

func (h *Httpx) Run(config *HttpxRunConfig) (chan string, error) {
	args := append(parseConfig(*config), "-silent", "-no-color")
	execPath, err := h.ExecPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe failed,err:", err)
		return nil, err
	}
	//获取标准输入
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("cmd.StdinPipe failed,err:", err)
		return nil, err
	}
	if config.URLCh != nil {
		//写入标准输入
		go func() {
			for domain := range config.URLCh {
				fmt.Fprintln(stdin, domain)
			}
			stdin.Close()
		}()
	}
	//fork子进程
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:" + err.Error())
		return nil, err
	}
	logger.Debugf("%s 启动成功\n", h.name)
	//读取标准输出
	resultCh := make(chan string, 1024)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			resultCh <- scanner.Text()
		}
		close(resultCh)
	}()
	return resultCh, err
}
