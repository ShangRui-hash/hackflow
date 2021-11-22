package hackflow

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/serkanalgur/phpfuncs"
)

type KSubdomain struct {
	name     string
	execPath string
}

//Name 返回工具名
func (k *KSubdomain) Name() string {
	return k.name
}

//ExecPath 返回工具执行路径
func (k *KSubdomain) ExecPath() (string, error) {
	if k.execPath == "" {
		if err := k.download(); err != nil {
			logger.Errorf("download %s failed,err:%v", k.Name(), err)
			return "", err
		}
	}
	return k.execPath, nil
}

//download 自动安装工具
func (k *KSubdomain) download() error {
	k.execPath = SavePath + "/ksubdomain/cmd/ksubdomain"
	dirPath := SavePath + "/ksubdomain"
	if !phpfuncs.FileExists(dirPath) {
		err := GetGit().Clone(CloneConfig{
			Url:      "https://github.com.cnpmjs.org/knownsec/ksubdomain",
			Depth:    1,
			SavePath: SavePath + "/ksubdomain/",
		})
		if err != nil {
			logger.Error("get clone failed,err:", err)
			return err
		}
	}
	if !phpfuncs.FileExists(k.execPath) {
		if err := GetGo().Mod(SavePath+"/ksubdomain", "download"); err != nil {
			logger.Error("get mod failed,err:", err)
			return err
		}
		err := GetGo().Build(BuildConfig{
			Path:  SavePath + "/ksubdomain/cmd",
			Files: []string{"ksubdomain.go"},
		})
		if err != nil {
			logger.Error("get build failed,err:", err)
		}
	}
	return nil
}

//GetKSubdomain 获取工具对象
func GetKSubdomain() *KSubdomain {
	if tool := container.Get(KSUBDOMAIN); tool != nil {
		return tool.(*KSubdomain)
	}
	container.Set(&KSubdomain{
		name: KSUBDOMAIN,
	})
	return container.Get(KSUBDOMAIN).(*KSubdomain)
}

type KSubdomainRunConfig struct {
	BruteLayer int    `flag:"-l"`
	Full       bool   `flag:"-full"`
	Verify     bool   `flag:"-verify"`
	Domain     string `flag:"-d"`
	DomainFile string `flag:"-dl'`
	DomainCh   chan string
}

func (k *KSubdomain) Run(config *KSubdomainRunConfig) (subdomainCh chan string, err error) {
	//构造命令
	args := append([]string{"-silent"}, parseConfig(*config)...)
	execPath, err := k.ExecPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, args...)
	//获取标准输出、标准错误输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe failed,err:", err)
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("cmd.StderrPipe failed,err:", err)
		return nil, err
	}
	output := io.MultiReader(stdout, stderr)
	//获取标准输入
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("cmd.StdinPipe failed,err:", err)
		return nil, err
	}
	if config.DomainCh != nil {
		//写入标准输入
		go func() {
			for domain := range config.DomainCh {
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
	logger.Debugf("%s 启动成功\n", k.name)
	//读取标准输出
	subdomainCh = make(chan string, 1024)
	go func() {
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			subdomainCh <- scanner.Text()
		}
		close(subdomainCh)
	}()

	return subdomainCh, err
}
