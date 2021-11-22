package main

import (
	"fmt"

	"github.com/ShangRui-hash/hackflow"
	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
	subfinder, err := hackflow.GetSubfinder()
	if err != nil {
		logrus.Error("get subfinder failed,err:", err)
		return
	}
	domainCh := make(chan string, 1024)
	domainList := []string{
		"lenovo.com",
		"lenovo.com.cn",
		"lenovomm.com",
		"lenovo.cn",
		"lenovo.net",
		"motorola.com",
		"motorola.com.cn",
		"baiying.cn",
	}
	go func() {
		for _, domain := range domainList {
			domainCh <- domain
		}
		close(domainCh)
	}()
	subdomainCh, err := subfinder.Run(hackflow.SubfinderRunConfig{
		Proxy:    "socks://127.0.0.1:7890",
		DomainCh: domainCh,
	})
	for subdomain := range subdomainCh {
		fmt.Println(subdomain)
	}
}
