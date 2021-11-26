package hackflow

import (
	"bufio"
	"os/exec"
	"regexp"
	"strings"

	"github.com/serkanalgur/phpfuncs"
)

type DirSearch struct {
	BaseTool
}

func newDirSearch() Tool {
	return &DirSearch{
		BaseTool: BaseTool{
			name: DIRSEARCH,
			desp: "目录扫描工具",
		},
	}
}

func GetDirSearch(isDebug bool) *DirSearch {
	return container.Get(DIRSEARCH).(*DirSearch)
}

type DirSearchResult struct {
	URL         string
	RespSize    string
	RespCode    string
	RedirectURL string
}

var resultReg = regexp.MustCompile(`^\[\d{2}:\d{2}:\d{2}\]\s{1}(\d{3})\s{1}-\s{2}([0-9A-Z]{4,5})\s+-\s+(.+)`)

//ExecPath 返回执行路径
func (d *DirSearch) ExecPath() (string, error) {
	return d.BaseTool.ExecPath(d.Download)
}

//download 下载源码
func (d *DirSearch) Download() (string, error) {
	execPath := SavePath + "/dirsearch/dirsearch.py"
	if !phpfuncs.FileExists(d.execPath) {
		logger.Debug("正在下载dirsearch")
		output, err := exec.Command("pip3", "install", "dirsearch", "--target="+SavePath).CombinedOutput()
		if err != nil {
			logger.Error("pip3 install failed,err:", err, "output:", string(output))
			return "", err
		}
		logger.Debug("下载源码完成，开始下载依赖:", string(output))
		output, err = exec.Command("pip3", "install", "-r", SavePath+"/dirsearch/requirements.txt").CombinedOutput()
		if err != nil {
			logger.Error("pip3 install -r failed,err:", err, "output:", string(output))
			return "", err
		}
		logger.Debug("依赖下载完成:", string(output))
	}
	return execPath, nil
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
