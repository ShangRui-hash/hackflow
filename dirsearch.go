package hackflow

import (
	"bufio"
	"os/exec"
	"regexp"
	"strings"

	"github.com/serkanalgur/phpfuncs"
	"github.com/sirupsen/logrus"
)

type DirSearch struct {
	name     string
	execPath string
	desp     string
}

func newDirSearch() Tool {
	return &DirSearch{
		name: DIRSEARCH,
		desp: "目录扫描工具",
	}
}

type DirSearchResult struct {
	URL         string
	RespSize    string
	RespCode    string
	RedirectURL string
}

var resultReg = regexp.MustCompile(`^\[\d{2}:\d{2}:\d{2}\]\s{1}(\d{3})\s{1}-\s{2}([0-9A-Z]{4,5})\s+-\s+(.+)`)

//Name 获取名称
func (d *DirSearch) Name() string {
	return d.name
}

//Desc 获取描述
func (d *DirSearch) Desp() string {
	return d.desp
}

//ExecPath 返回执行路径
func (d *DirSearch) ExecPath() (string, error) {
	if len(d.execPath) == 0 {
		if err := d.download(); err != nil {
			logger.Errorf("download %s failed,err:", err)
			return "", err
		}
	}
	return d.execPath, nil
}

//download 下载源码
func (d *DirSearch) download() error {
	d.execPath = SavePath + "/dirsearch/dirsearch.py"
	if !phpfuncs.FileExists(d.execPath) {
		logger.Debug("正在下载dirsearch")
		output, err := exec.Command("pip3", "install", "dirsearch", "--target="+SavePath).CombinedOutput()
		if err != nil {
			logger.Error("pip3 install failed,err:", err, "output:", string(output))
			return err
		}
		logger.Debug("下载源码完成，开始下载依赖:", string(output))
		output, err = exec.Command("pip3", "install", "-r", SavePath+"/dirsearch/requirements.txt").CombinedOutput()
		if err != nil {
			logger.Error("pip3 install -r failed,err:", err, "output:", string(output))
			return err
		}
		logger.Debug("依赖下载完成:", string(output))
	}
	return nil
}

func GetDirSearch(isDebug bool) (*DirSearch, error) {
	if tool := container.Get(DIRSEARCH); tool != nil {
		return tool.(*DirSearch), nil
	}
	tool := &DirSearch{
		name: DIRSEARCH,
	}
	if err := tool.download(); err != nil {
		logrus.Error("tool.download failed,err:", err)
		return nil, err
	}
	container.Set(tool)
	return tool, nil
}

type DirSearchConfig struct {
	TargetURL       string   `flag:"-u"`
	HTTPMethod      string   `flag:"-m"`
	FullURL         bool     `flag:"--full-url"`
	RandomAgent     bool     `flag:"--random-agent"`
	RemoveExtension bool     `flag:"--remove-extensions"`
	EXT             []string `flag:"-e"`
	Subdirs         []string `flag:"--subdirs"`
}

func (d *DirSearch) Run(config DirSearchConfig) (chan *DirSearchResult, error) {
	execPath, err := d.ExecPath()
	if err != nil {
		return nil, err
	}
	args := []string{execPath, "--no-color"}
	args = append(args, parseConfig(config)...)
	cmd := exec.Command("python3", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe failed,err:", err)
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:" + err.Error())
		return nil, err
	}
	urlCh := make(chan *DirSearchResult, 1024)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if result := d.ParseResult(scanner.Text()); result != nil {
				urlCh <- result
			}
		}
		close(urlCh)
	}()
	return urlCh, nil
}

func (d *DirSearch) ParseResult(line string) (result *DirSearchResult) {
	if resultReg.MatchString(line) {
		submatch := resultReg.FindStringSubmatch(line)
		result = &DirSearchResult{
			RespSize: submatch[2],
			RespCode: submatch[1],
		}
		if strings.Contains(submatch[3], "->") {
			urls := strings.Split(submatch[3], "->")
			result.URL = urls[0]
			result.RedirectURL = urls[1]
		} else {
			result.URL = submatch[3]
		}
	}
	return result
}
