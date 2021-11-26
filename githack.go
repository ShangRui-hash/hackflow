package hackflow

import (
	"path/filepath"
)

type GitHack struct {
	BaseTool
}

func newGitHack() Tool {
	return &GitHack{
		BaseTool: BaseTool{
			name: GIT_HACK,
			desp: "Git 泄漏利用工具",
		},
	}
}

func GetGitHack() *GitHack {
	return container.Get(GIT_HACK).(*GitHack)
}

func (g *GitHack) Download() (string, error) {
	err := GetGit().Clone(CloneConfig{
		Url:      "https://github.com.cnpmjs.org/lijiejie/GitHack",
		SavePath: filepath.Join(SavePath, g.name),
		Depth:    1,
	})
	if err != nil {
		logger.Error("git clone failed,err:", err)
		return "", err
	}
	return filepath.Join(SavePath, g.name, "GitHack.py"), nil
}

func (g *GitHack) ExecPath() (string, error) {
	return g.BaseTool.ExecPath(g.Download)
}
