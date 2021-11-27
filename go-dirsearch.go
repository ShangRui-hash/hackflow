package hackflow

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/evilsocket/brutemachine"
	"github.com/serkanalgur/phpfuncs"
)

type GoDirsearchRunConfig struct {
	TargetCh            chan string
	RoutineCount        int
	RandomAgent         bool
	Proxy               string
	Dictionary          string
	StatusCodeBlackList []int
}

func GoDirsearch(config *GoDirsearchRunConfig) (chan string, error) {
	resultCh := make(chan string, 1024)
	for i := 0; i < config.RoutineCount; i++ {
		go func() {
			for targetURL := range config.TargetCh {
				m := brutemachine.New(-1, config.Dictionary, func(page string) interface{} {
					url := fmt.Sprintf("%s/%s", targetURL, page)
					resp, err := http.Head(url)
					if err != nil {
						logger.Error("http.Head failed,err:", err)
						return nil
					}
					if !phpfuncs.InArray(resp.StatusCode, config.StatusCodeBlackList) {
						return strings.Join([]string{url, fmt.Sprintf("%d", resp.StatusCode)}, ":")
					}
					return nil
				}, func(res interface{}) {
					fmt.Printf("@ Found '%s'\n", res)
					resultCh <- res.(string)
				})
				if err := m.Start(); err != nil {
					panic(err)
				}
				m.Wait()
			}
		}()
	}
	return resultCh, nil
}
