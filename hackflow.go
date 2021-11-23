package hackflow

import (
	"fmt"
	"go/build"
	"reflect"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	GO            = "go"
	GIT           = "git"
	SQLMAP        = "sqlmap"
	URL_COLLECTOR = "url_collector"
	DIRSEARCH     = "dirsearch"
	KSUBDOMAIN    = "ksubdomain"
	SUBFINDER     = "subfinder"
	HTTPX         = "httpx"
)

var (
	container *Container
	logger    = logrus.New()
	SavePath  = build.Default.GOPATH + "/hackflow"
	newTool   = map[string]func() Tool{
		GO:            newGo,
		GIT:           newGit,
		SQLMAP:        newSqlmap,
		URL_COLLECTOR: newUrlCollector,
		DIRSEARCH:     newDirSearch,
		KSUBDOMAIN:    newKSubdomain,
		SUBFINDER:     newSubfinder,
		HTTPX:         newHttpx,
	}
)

type Container struct {
	allTools []string
	tools    sync.Map
}

type Tool interface {
	Name() string
	Desp() string
	ExecPath() (string, error)
	download() error
}

func init() {
	container = &Container{
		tools: sync.Map{},
		allTools: []string{
			SQLMAP,
			URL_COLLECTOR,
			DIRSEARCH,
			KSUBDOMAIN,
			SUBFINDER,
			HTTPX,
		},
	}
}

//Set 将工具注册到注册树
func (c *Container) Set(tool Tool) {
	c.tools.Store(tool.Name(), tool)
}

//Get 从注册树上获取工具
func (c *Container) Get(name string) Tool {
	if tool, ok := c.tools.Load(name); ok {
		return tool.(Tool)
	}
	tool := newTool[name]()
	container.Set(tool)
	return tool
}

//GetAllTools 获取全部工具
func GetAllTools() (tools []Tool) {
	for i := range container.allTools {
		tools = append(tools, container.Get(container.allTools[i]))
	}
	return tools
}

//SetDebug 是否开启debug模式
func SetDebug(debug bool) {
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}

//parseConfig 解析config
func parseConfig(config interface{}) (args []string) {
	typeof := reflect.TypeOf(config)
	valueof := reflect.ValueOf(config)
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		switch valueof.Field(i).Type().Kind() {
		case reflect.Bool:
			if valueof.Field(i).Bool() {
				args = append(args, field.Tag.Get("flag"))
			}
		case reflect.Int:
			if valueof.Field(i).Int() != 0 {
				args = append(args, []string{field.Tag.Get("flag"), fmt.Sprintf("%d", valueof.Field(i).Int())}...)
			}
		case reflect.String:
			if valueof.Field(i).String() != "" {
				args = append(args, []string{field.Tag.Get("flag"), valueof.Field(i).String()}...)
			}
		}
	}
	logger.Debug("parse config to args:", args)
	return args
}
