package main

import (
	"fmt"

	"github.com/ShangRui-hash/hackchain"
	"github.com/sirupsen/logrus"
)

func main() {
	// URLCollector, err := hackchain.GetUrlCollector()
	// if err != nil {
	// 	logrus.Error("hackchain.GetUrlCollector failed,err:", err)
	// 	return
	// }
	// urlCh, err := URLCollector.Run(&hackchain.UrlCollectorCofnig{
	// 	Keyword:      ".php?id=10",
	// 	SearchEngine: "google",
	// 	Proxy:        "http://127.0.0.1:7890",
	// })
	// if err != nil {
	// 	logrus.Error("URLCollector.Run failed,err:", err)
	// 	return
	// }
	// for url := range urlCh {
	// 	fmt.Println(url)
	// }

	dirSearch, err := hackchain.GetDirSearch(true)
	if err != nil {
		logrus.Error("hackchain.GetDirSearch failed,err:", err)
		return
	}
	dirSearch.SetDebug(false)
	resultCh, err := dirSearch.Run(hackchain.DirSearchConfig{
		TargetURL: "www.baidu.com",
		FullURL:   true,
	})
	if err != nil {
		logrus.Error("dirSearch.Run failed,err:", err)
		return
	}
	for result := range resultCh {
		fmt.Printf("%+v\n", result)
	}

}
